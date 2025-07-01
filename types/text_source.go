package types

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Mode string

type TextSource struct {
	URL    string
	Parser func([]byte) (string, error)
}

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
