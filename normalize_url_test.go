package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove https://",
			inputURL: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove http://",
			inputURL: "http://www.youtube.com/",
			expected: "www.youtube.com",
		},
		{
			name:     "retain subdomain",
			inputURL: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "handle uppercase",
			inputURL: "HTTPS://www.youtube.com/",
			expected: "www.youtube.com",
		},
		{
			name:     "Ignore fragments",
			inputURL: "https://blog.boot.dev/path#intro",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove trailing slash",
			inputURL: "https://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "root path with trailing slash",
			inputURL: "https://example.com/",
			expected: "example.com",
		},
		{
			name:     "root path without slash",
			inputURL: "https://example.com",
			expected: "example.com",
		},
		{
			name:     "port number",
			inputURL: "http://localhost:8080/docs/",
			expected: "localhost:8080/docs",
		},
		{
			name:     "ignore query parameters",
			inputURL: "https://boot.dev/search?q=golang",
			expected: "boot.dev/search",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v -  %s. Unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s. Expected URL %v, actual URL %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
