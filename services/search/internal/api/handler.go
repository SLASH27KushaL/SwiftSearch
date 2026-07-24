package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"moogle-go/services/search/internal/engine"
	"moogle-go/services/search/internal/store"
	"moogle-go/services/search/pkg/models"
)

type SearchHandler struct {
	ranker *engine.Ranker
	cache  *store.RedisCache
}

func NewSearchHandler(ranker *engine.Ranker, cache *store.RedisCache) *SearchHandler {
	return &SearchHandler{
		ranker: ranker,
		cache:  cache,
	}
}

func (h *SearchHandler) HandleSearch(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'q' parameter"})
		return
	}

	startTime := time.Now()
	cacheKey := "search:" + query

	// 1. Check Redis Cache First (O(1) retrieval)
	if cachedData, err := h.cache.Get(c.Request.Context(), cacheKey); err == nil {
		var response models.SearchResponse
		json.Unmarshal([]byte(cachedData), &response)
		response.ExecutionTimeMs = time.Since(startTime).Milliseconds() // Update execution time
		c.JSON(http.StatusOK, response)
		return
	}

	// 2. Cache Miss -> Execute Ranking Engine
	results, err := h.ranker.ExecuteSearch(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal search error"})
		return
	}

	// 3. Construct JSON Payload
	response := models.SearchResponse{
		Query:           query,
		TotalResults:    len(results),
		ExecutionTimeMs: time.Since(startTime).Milliseconds(),
		Results:         results,
	}

	// 4. Save to Cache asynchronously (TTL: 5 minutes)
	go func(resp models.SearchResponse) {
		jsonData, _ := json.Marshal(resp)
		_ = h.cache.Set(context.Background(), cacheKey, string(jsonData), 5*time.Minute)
	}(response)

	// 5. Send back to Client
	c.JSON(http.StatusOK, response)
}