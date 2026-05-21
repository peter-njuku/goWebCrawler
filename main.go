package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

func main() {
	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Println("No website provided")
		os.Exit(1)
	}

	if len(arguments) > 2 {
		fmt.Println("Too many arguments provided")
		os.Exit(1)
	}

	if len(arguments) == 2 {
		baseURL := arguments[1]
		fmt.Println("Starting crawl of:", baseURL)
		html, err := getHTML(baseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(html)
	}
}
