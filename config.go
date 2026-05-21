package main

import (
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, ok := cfg.pages[normalizedURL]; ok {
		// data.visitCount++
		// cfg.pages[normalizedURL] = data
		return false
	}
	if len(cfg.pages) >= cfg.maxPages {
		return false
	}
	cfg.pages[normalizedURL] = PageData{URL: normalizedURL, visitCount: 1}
	return true
}

func (cfg *config) setPageData(normalizedURL string, data PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	data.visitCount = cfg.pages[normalizedURL].visitCount
	cfg.pages[normalizedURL] = data
}

func configure(rawBaseURL string, maxConcurrency, maxPages int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	return &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}, nil
}

func (cfg *config) pagesLen() int {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages)
}
