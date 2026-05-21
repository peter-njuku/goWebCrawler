package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return ""
	}

	h1 := doc.Find("h1").First().Text()
	if trimmed := strings.TrimSpace(h1); trimmed != "" {
		return trimmed
	}

	h2 := doc.Find("h2").First().Text()
	return strings.TrimSpace(h2)
}

func getFirstParagraphFromHTML(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return ""
	}

	main := doc.Find("main")
	if main.Length() > 0 {
		p := main.Find("p").First().Text()
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			return trimmed
		}
	}

	p := doc.Find("p").First().Text()
	return strings.TrimSpace(p)
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	urls := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		parsedURL, err := url.Parse(href)
		if err != nil {
			return
		}

		resolvedURL := baseURL.ResolveReference(parsedURL)
		urls = append(urls, resolvedURL.String())
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	images := []string{}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		image, exists := s.Attr("src")
		if !exists {
			return
		}

		parsedURL, err := url.Parse(image)
		if err != nil {
			return
		}

		resolvedURL := baseURL.ResolveReference(parsedURL)
		images = append(images, resolvedURL.String())
	})
	return images, nil
}
