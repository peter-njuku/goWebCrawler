# Go Web Scraper

A lightweight Go-based web scraping starter project for extracting page metadata from HTML.

## Features

- Parses HTML to extract the first heading and first paragraph.
- Resolves outgoing links from `<a>` tags against the page base URL.
- Collects image source URLs from `<img>` tags.
- Normalizes URLs for consistent comparison and deduplication.
- Includes unit tests for HTML parsing and URL normalization.

## Requirements

- Go 1.25 or newer

## Installation

```bash
go mod tidy
```

## Usage

Run the application with a single target URL:

```bash
go run . https://example.com
```

The CLI validates the argument count and prints the starting crawl URL.

## Testing

Run the unit test suite with:

```bash
go test ./...
```

## Project structure

- `main.go` - command-line entrypoint and argument validation
- `parseHTML.go` - HTML extraction helpers for headings, paragraphs, links, and images
- `normalize_url.go` - URL normalization logic
- `normalize_url_test.go` - URL normalization tests
- `parseHTML_test.go` - HTML parsing and page extraction tests

## Notes

This project is a solid foundation for building a larger web scraper or crawler. The current implementation focuses on parsing and normalization utilities that can be extended with fetching, queuing, and persistence logic.
