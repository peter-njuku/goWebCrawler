package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	path := parsedURL.Path
	if path != "" && path != "/" {
		path = strings.TrimSuffix(path, "/")
	} else if path == "/" {
		path = ""
	}
	return parsedURL.Host + path, nil
}
