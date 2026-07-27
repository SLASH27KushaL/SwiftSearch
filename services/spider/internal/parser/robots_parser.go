package parser

import (
	"net/http"
	"net/url"
	"time"

	"github.com/temoto/robotstxt"
)

// RobotsChecker manages fetching and evaluating robots.txt rules
type RobotsChecker struct {
	userAgent string
	client    *http.Client
}

// NewRobotsChecker initializes the checker with your bot's custom name
func NewRobotsChecker(userAgent string) *RobotsChecker {
	return &RobotsChecker{
		userAgent: userAgent,
		client: &http.Client{
			Timeout: 5 * time.Second, // Don't hang forever waiting for a robots.txt
		},
	}
}

// IsAllowed checks if the target URL is legally allowed to be crawled
func (rc *RobotsChecker) IsAllowed(targetURL string) bool {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return false // If the URL is broken, don't crawl it
	}

	// Construct the URL to the domain's robots.txt
	robotsURL := parsedURL.Scheme + "://" + parsedURL.Host + "/robots.txt"

	resp, err := rc.client.Get(robotsURL)
	if err != nil {
		// If the server is unreachable or robots.txt doesn't exist,
		// standard web etiquette says it's safe to crawl.
		return true
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// 404 means no robots.txt exists. Crawl away!
		return true
	}

	robotsData, err := robotstxt.FromResponse(resp)
	if err != nil {
		return true // If they wrote a malformed robots file, we assume allowed
	}

	// Find the rules for our specific bot (or the default "*" rules)
	group := robotsData.FindGroup(rc.userAgent)

	// Test if the specific path (e.g., "/private/data") is allowed
	return group.Test(parsedURL.Path)
}
