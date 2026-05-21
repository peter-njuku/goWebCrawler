package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "simple h1",
			html:     "<html><body><h1>Hello</h1></body></html>",
			expected: "Hello",
		},
		{
			name:     "fallback to h2",
			html:     "<h2>Subtitle</h2><p>ignore</p>",
			expected: "Subtitle",
		},
		{
			name:     "no heading",
			html:     "<p>just text</p>",
			expected: "",
		},
		{
			name:     "nested tags inside h1",
			html:     "<h1>Welcome to <strong>Boot.dev</strong></h1>",
			expected: "Welcome to Boot.dev",
		},
		{
			name:     "h1 after h2, should return h1",
			html:     "<h2>ignore</h2><h1>main title</h1>",
			expected: "main title",
		},
		{
			name:     "empty h1",
			html:     "<h1></h1><h2>fallback</h2>",
			expected: "fallback",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := getHeadingFromHTML(tc.html)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "paragraph inside main",
			html:     "<main><p>Inside main</p></main><p>Outside</p>",
			expected: "Inside main",
		},
		{
			name:     "no main, first paragraph",
			html:     "<div><p>First p</p><p>Second</p></div>",
			expected: "First p",
		},
		{
			name:     "no paragraphs",
			html:     "<h1>Hello</h1>",
			expected: "",
		},
		{
			name:     "main exists but no p inside",
			html:     "<main><h1>Title</h1></main><p>Outside</p>",
			expected: "Outside",
		},
		{
			name:     "nested tags inside p",
			html:     "<main><p>Learn <strong>coding</strong> now</p></main>",
			expected: "Learn coding now",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := getFirstParagraphFromHTML(tc.html)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	type testCase struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
		expectErr bool
	}

	tests := []testCase{
		{
			name:      "absolute URL",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com"},
			expectErr: false,
		},
		{
			name:      "relative URL with no path",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="/some-path"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com/some-path"},
			expectErr: false,
		},
		{
			name:      "relative URL with dot segments",
			inputURL:  "https://crawler-test.com/docs/",
			inputBody: `<html><body><a href="./api"><span>API</span></a><a href="../home">Home</a></body></html>`,
			expected:  []string{"https://crawler-test.com/docs/api", "https://crawler-test.com/home"},
			expectErr: false,
		},
		{
			name:     "multiple anchors with mix of absolute and relative",
			inputURL: "https://crawler-test.com/base/",
			inputBody: `<html><body>
                <a href="https://crawler-test.com/absolute">Absolute</a>
                <a href="relative">Relative</a>
                <a href="/root-relative">Root Relative</a>
                <a href="./current-dir">Current Dir</a>
                <a href="../parent">Parent</a>
            </body></html>`,
			expected: []string{
				"https://crawler-test.com/absolute",
				"https://crawler-test.com/base/relative",
				"https://crawler-test.com/root-relative",
				"https://crawler-test.com/base/current-dir",
				"https://crawler-test.com/parent",
			},
			expectErr: false,
		},
		{
			name:      "anchor with no href (ignore)",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a name="top">Just a named anchor</a></body></html>`,
			expected:  []string{},
			expectErr: false,
		},
		{
			name:      "anchor with empty href",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="">Empty link</a></body></html>`,
			expected:  []string{"https://crawler-test.com"}, // or empty? Usually resolves to base URL
			expectErr: false,
		},
		{
			name:      "href with fragment only",
			inputURL:  "https://crawler-test.com/page",
			inputBody: `<html><body><a href="#section">Jump to section</a></body></html>`,
			expected:  []string{"https://crawler-test.com/page#section"},
			expectErr: false,
		},
		{
			name:      "href with query parameters",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="/search?q=golang&page=2">Search</a></body></html>`,
			expected:  []string{"https://crawler-test.com/search?q=golang&page=2"},
			expectErr: false,
		},
		{
			name:      "malformed HTML (missing closing tag)",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="/broken"><span>Link</span>`,
			expected:  []string{"https://crawler-test.com/broken"},
			expectErr: false,
		},
		{
			name:      "invalid base URL",
			inputURL:  "://invalid",
			inputBody: `<html><body><a href="/test">Test</a></body></html>`,
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "no anchor tags at all",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><p>No links here</p></body></html>`,
			expected:  []string{},
			expectErr: false,
		},
		{
			name:      "href uses javascript: (should ignore or treat as is?)",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="javascript:void(0)">Click</a></body></html>`,
			expected:  []string{"javascript:void(0)"}, // depending on design; often filtered out
			expectErr: false,
		},
		{
			name:      "href with mailto:",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="mailto:test@example.com">Email</a></body></html>`,
			expected:  []string{"mailto:test@example.com"}, // maybe keep or skip
			expectErr: false,
		},
		{
			name:     "duplicate URLs",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body>
                <a href="/dupe">Link1</a>
                <a href="https://crawler-test.com/dupe">Link2</a>
            </body></html>`,
			expected:  []string{"https://crawler-test.com/dupe", "https://crawler-test.com/dupe"}, // or deduplicated
			expectErr: false,
		},
		{
			name:      "nested HTML elements",
			inputURL:  "https://crawler-test.com",
			inputBody: `<div><a href="/nested"><span><strong>Deep</strong></span></a></div>`,
			expected:  []string{"https://crawler-test.com/nested"},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseURL, err := url.Parse(tt.inputURL)
			if err != nil {
				if !tt.expectErr {
					t.Fatalf("Failed to parse base URL: %v", err)
				}
				return
			}

			actual, err := getURLsFromHTML(tt.inputBody, baseURL)
			if tt.expectErr && err == nil {
				t.Fatalf("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.expectErr {
				return
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Expected %v : got %v", tt.expected, actual)
			}
		})
	}
}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><img src="/logo.png" alt="Logo"></body></html>`

    baseURL, err := url.Parse(inputURL)
    if err != nil {
        t.Errorf("couldn't parse input URL: %v", err)
        return
    }

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}