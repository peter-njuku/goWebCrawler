# Go Web Crawler

A lightweight Go web crawler that fetches HTML pages, extracts metadata, resolves internal links, and writes a JSON report.

## Features

- Fetches pages over HTTP with a custom `User-Agent`.
- Validates HTML content by checking the response `Content-Type`.
- Extracts the first heading (`<h1>` or fallback to `<h2>`).
- Extracts the first paragraph, preferring content inside `<main>`.
- Resolves outgoing links from `<a href="...">` tags relative to the page base URL.
- Extracts image source URLs from `<img src="...">` tags.
- Normalizes URLs for consistent deduplication and page tracking.
- Crawls within the same hostname only.
- Uses concurrency controls and a maximum page limit.
- Writes a sorted JSON report to `report.json`.

## Requirements

- Go 1.25 or newer

## Installation

From the project root:

```bash
go mod tidy
```

## Usage

Run the crawler with three arguments:

```bash
go run . <start-url> <max-concurrency> <max-pages>
```

Example:

```bash
go run . https://example.com 4 50
```

The application will:

- validate the provided arguments,
- start crawling from the provided URL,
- follow outgoing links only on the same hostname,
- stop after reaching the maximum number of pages,
- and write the crawl results to `report.json`.

## Output

The crawler produces a JSON report named `report.json` containing an array of pages with:

- `url`
- `heading`
- `first_paragraph`
- `outgoing_links`
- `image_urls`

The report is written with pages sorted by normalized URL.

## Project structure

- `main.go` - CLI entrypoint, argument validation, crawl orchestration, report generation
- `config.go` - crawler configuration, concurrency control, page tracking, and synchronization
- `parseHTML.go` - HTML parsing helpers and page metadata extraction
- `normalize_url.go` - URL normalization logic for deduplication
- `json_report.go` - JSON report generation and file writing
- `normalize_url_test.go` - unit tests for URL normalization
- `parseHTML_test.go` - unit tests for HTML parsing and page extraction

## Testing

Run the full unit test suite:

```bash
go test ./...
```

## Notes

- The crawler currently does not persist state beyond the generated JSON report.
- It only crawls pages under the same hostname as the start URL.
- The extracted page metadata is intended as a foundation for more advanced crawling, indexing, or scraping features.
