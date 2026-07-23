package crawler

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"moogle-go/services/spider/internal/parser"
	"moogle-go/services/spider/internal/storage"
	"moogle-go/services/spider/pkg/models"
)

// EngineOptions holds the configuration for the crawler engine.
type EngineOptions struct {
	WorkerCount  int
	MaxDepth     int
	CrawlDelay   time.Duration
	UserAgent    string
	VisitedStore *storage.RedisVisitedStore
	DocStore     *storage.MongoStore
	Parser       *parser.HTMLParser
	Normalizer   *parser.URLNormalizer
}

// Engine orchestrates the crawling process, managing the job queue and worker pool.
type Engine struct {
	opts     EngineOptions
	jobQueue chan models.CrawlJob
	wg       sync.WaitGroup
	quit     chan struct{}
	client   *http.Client
}

// NewEngine initializes a new Engine instance.
func NewEngine(opts EngineOptions) *Engine {
	return &Engine{
		opts:     opts,
		// Buffer the job queue to prevent blocking when submitting initial seeds
		jobQueue: make(chan models.CrawlJob, 100000),
		quit:     make(chan struct{}),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Start launches the worker pool and blocks until stopped.
func (e *Engine) Start(ctx context.Context) {
	log.Printf("Starting %d crawler workers...", e.opts.WorkerCount)
	for i := 0; i < e.opts.WorkerCount; i++ {
		e.wg.Add(1)
		go e.worker(ctx, i)
	}
	e.wg.Wait()
}

// Stop initiates a graceful shutdown of all workers.
func (e *Engine) Stop() {
	close(e.quit)
}

// SubmitJob adds a new URL to the crawl queue if it hasn't exceeded the max depth.
func (e *Engine) SubmitJob(job models.CrawlJob) {
	if job.Depth > e.opts.MaxDepth {
		return
	}

	// Non-blocking send to avoid deadlocks if the queue is full
	select {
	case e.jobQueue <- job:
	default:
		log.Printf("Job queue is full, dropping URL: %s", job.URL)
	}
}

// worker represents a single concurrent goroutine fetching and parsing pages.
func (e *Engine) worker(ctx context.Context, id int) {
	defer e.wg.Done()
	for {
		select {
		case <-ctx.Done(): // Context cancelled (e.g., Ctrl+C)
			return
		case <-e.quit: // Engine stopped gracefully
			return
		case job := <-e.jobQueue:
			e.processJob(ctx, job, id)
			time.Sleep(e.opts.CrawlDelay) // Politeness delay
		}
	}
}

// processJob handles the entire lifecycle of a single URL crawl.
func (e *Engine) processJob(ctx context.Context, job models.CrawlJob, workerID int) {
	// 1. Deduplication check
	visited, err := e.opts.VisitedStore.IsVisited(ctx, job.URL)
	if err != nil || visited {
		return
	}

	// 2. Mark as visited immediately to prevent race conditions
	_ = e.opts.VisitedStore.MarkVisited(ctx, job.URL)

	// 3. HTTP Fetch
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.URL, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", e.opts.UserAgent)

	resp, err := e.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Only process successful HTML responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// 4. Parse DOM
	page, err := e.opts.Parser.Parse(job.URL, resp.StatusCode, body)
	if err != nil {
		return
	}

	// 5. Persist to MongoDB
	err = e.opts.DocStore.SavePage(ctx, page)
	if err != nil {
		log.Printf("[Worker %d] Failed to save page %s: %v", workerID, job.URL, err)
		return
	}

	log.Printf("[Worker %d] Crawled: %s (Found %d links)", workerID, job.URL, len(page.Outlinks))

	// 6. Queue newly discovered links
	if job.Depth < e.opts.MaxDepth {
		for _, link := range page.Outlinks {
			e.SubmitJob(models.CrawlJob{
				URL:   link.URL,
				Depth: job.Depth + 1,
			})
		}
	}
}