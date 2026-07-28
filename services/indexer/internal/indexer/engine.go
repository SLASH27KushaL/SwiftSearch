package indexer

import (
	"context"
	"log"

	"moogle-go/services/indexer/internal/store"
	"moogle-go/services/indexer/internal/tokenizer"
	"moogle-go/services/indexer/pkg/models"
)

type Engine struct {
	reader    *store.MongoReader
	writer    *store.IndexWriter
	batchSize int
}

func NewEngine(reader *store.MongoReader, writer *store.IndexWriter, batchSize int) *Engine {
	return &Engine{
		reader:    reader,
		writer:    writer,
		batchSize: batchSize,
	}
}

func (e *Engine) Run(ctx context.Context) error {
	for {
		pages, err := e.reader.GetNextBatch(ctx, e.batchSize)
		if err != nil {
			return err
		}

		if len(pages) == 0 {
			log.Println("No more new documents to index. Pipeline resting...")
			break
		}

		log.Printf("Processing batch of %d documents...", len(pages))
		var indexBatch []models.IndexEntry

		for _, page := range pages {
			tokens := tokenizer.CleanAndTokenize(page.TextContent)
			filteredTokens := tokenizer.RemoveStopwords(tokens)
			tfMap := CalculateTF(filteredTokens)

			for term, tf := range tfMap {
				// FIX: We now properly use the 'Matches' array with our single match
				entry := models.IndexEntry{
					Term: term,
					Matches: []models.DocumentMatch{
						{
							URL:   page.URL,
							Title: page.Title,
							TF:    tf,
						},
					},
				}
				indexBatch = append(indexBatch, entry)
			}

			if err := e.reader.MarkAsProcessed(ctx, page.URL); err != nil {
				log.Printf("Failed to mark %s as processed: %v", page.URL, err)
			}
		}

		if len(indexBatch) > 0 {
			// Chunk the massive array into batches of 500 to prevent MongoDB EOF crashes
			chunkSize := 500
			for i := 0; i < len(indexBatch); i += chunkSize {
				end := i + chunkSize
				if end > len(indexBatch) {
					end = len(indexBatch)
				}

				chunk := indexBatch[i:end]
				if err := e.writer.BulkUpsert(ctx, chunk); err != nil {
					log.Printf("Failed to write index chunk to database: %v", err)
					return err
				}
			}
		}
	}

	return nil
}