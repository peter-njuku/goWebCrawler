package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP Error: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf("Invalid Content-Type: %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	if cfg.pagesLen() >= cfg.maxPages {
		return
	}

	currentParsed, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error parsing current URL: %v\n", err)
		return
	}

	if currentParsed.Hostname() != cfg.baseURL.Hostname() {
		return
	}

	normalized, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error normalizing URL: %v\n", err)
		return
	}

	//fmt.Printf("DEBUG: normalized = %s\n", normalized)

	// if count, exists := pages[normalized]; exists {
	// 	pages[normalized] = count + 1
	// 	//fmt.Printf("DEBUG: already seen %s, count now %d\n", normalized, pages[normalized])
	// 	return
	// }
	// pages[normalized] = 1
	// //fmt.Printf("DEBUG: new page %s\n", normalized)
	isFirst := cfg.addPageVisit(normalized)
	if !isFirst {
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error fetching HTML: %v\n", err)
		return
	}

	// urls, err := getURLsFromHTML(html, baseParsed)
	// if err != nil {
	// 	fmt.Printf("Error extracting URLs: %v\n", err)
	// 	return
	// }
	pageData := extractPageData(html, rawCurrentURL)
	cfg.setPageData(normalized, pageData)

	for _, nextURL := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(nextURL)
		// fmt.Println(nextURL)
		// crawlPage(rawBaseURL, nextURL, pages)
	}
}

func main() {
	arguments := os.Args

	if len(arguments) < 4 {
		fmt.Println("No website provided")
		os.Exit(1)
	}

	if len(arguments) > 4 {
		fmt.Println("Too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := arguments[1]
	maxConcurrencyString := os.Args[2]
	maxPagesString := os.Args[3]

	maxConcurrency, err := strconv.Atoi(maxConcurrencyString)
	if err != nil || maxConcurrency <= 0 {
		fmt.Printf("Invalid max concurrency: %s\n", maxConcurrencyString)
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(maxPagesString)
	if err != nil || maxPages <= 0 {
		fmt.Printf("Invalid max pages: %s\n", maxPagesString)
		os.Exit(1)
	}

	cfg, err := configure(rawBaseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("Error configuring crawler: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Starting crawl of:", rawBaseURL)

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()
	for normalizedURL := range cfg.pages {
		fmt.Printf("Found: %s\n", normalizedURL)
	}
	
}
