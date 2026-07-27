package worker

import (
	"fmt"
	"log"
	"sync"
	"time"

	"swiftsearch/spider/internal/filter"
	"swiftsearch/spider/internal/frontier"
	"swiftsearch/spider/internal/parser"
	"swiftsearch/spider/internal/store"
)

type Pool struct {
	Queue      *frontier.RedisQueue
	Filter     *filter.BloomFilter
	Store      *store.MongoWriter
	Politeness *frontier.PolitenessPolicy // Add this here
}

func NewPool(q *frontier.RedisQueue, f *filter.BloomFilter, s *store.MongoWriter, p *frontier.PolitenessPolicy) *Pool {
	return &Pool{Queue: q, Filter: f, Store: s, Politeness: p}
}

func (p *Pool) Start(workerCount int) {
	var wg sync.WaitGroup

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go p.worker(i, &wg)
	}

	log.Printf("Started %d distributed spider workers...", workerCount)
	wg.Wait()
}

func (p *Pool) worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// 1. Pop URL from Redis Queue
		rawURL, err := p.Queue.Pop()
		if err != nil {
			continue
		}

		// 2. Enforce Politeness
		domain := frontier.ExtractDomain(rawURL)
		if !p.Politeness.RequestAccess(domain) {
			// Domain was hit too recently by the cluster.
			// Push it back to the queue and take a tiny breath to prevent CPU thrashing.
			p.Queue.Push(rawURL)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		fmt.Printf("[Worker %d] Crawling: %s\n", id, rawURL)

		// 3. Fetch HTML and parse links
		page, err := parser.FetchAndParse(rawURL)
		if err != nil {
			continue
		}

		// 4. Save raw HTML to MongoDB
		err = p.Store.SavePage(page.URL, page.HTML)
		if err != nil {
			log.Println("Mongo Error:", err)
		}

		// 5. Check new links against Bloom Filter and queue unvisited ones
		for _, link := range page.Links {
			visited, _ := p.Filter.CheckAndAdd(link)
			if !visited {
				p.Queue.Push(link)
			}
		}
	}
}
