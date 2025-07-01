package types

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Mode string // Mode is a type that represents the different modes of text formatting.
// TextSource is a struct that represents a source of text.
// It contains a URL to fetch the text from and a parser function to process the response.
type TextSource struct {
	URL    string                       // URL to fetch text from
	Parser func([]byte) (string, error) // Parser function to process the response
}

// FetchText fetches text from the specified URL and parses it using the provided parser function.
func (s *TextSource) FetchText() (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Printf("Fetching from URL: %s\n", s.URL)
	resp, err := client.Get(s.URL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if s.Parser != nil {
		return s.Parser(body)
	}
	return string(body), nil
}
