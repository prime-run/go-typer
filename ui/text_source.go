package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type TextSource interface {
	FetchText() (string, error)
	FormatText(text string) string
}

type ZenQuotesSource struct {
	URL string
}

func NewZenQuotesSource() *ZenQuotesSource {
	return &ZenQuotesSource{
		URL: "https://zenquotes.io/api/random",
	}
}

func (s *ZenQuotesSource) FetchText() (string, error) {
	DebugLog("TextSource: Fetching quote from %s", s.URL)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(s.URL)
	if err != nil {
		DebugLog("TextSource: Failed to fetch quote: %v", err)
		return "", fmt.Errorf("failed to fetch quote: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		DebugLog("TextSource: API returned non-200 status: %d", resp.StatusCode)
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		DebugLog("TextSource: Failed to read response: %v", err)
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	DebugLog("TextSource: Raw API response: %s", string(body))

	var result []struct {
		Quote  string `json:"q"`
		Author string `json:"a"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		DebugLog("TextSource: Failed to parse quote: %v", err)
		return "", fmt.Errorf("failed to parse quote: %w", err)
	}

	if len(result) == 0 {
		DebugLog("TextSource: No quotes returned from API")
		return "", fmt.Errorf("no quotes returned from API")
	}

	quote := result[0].Quote
	author := result[0].Author

	// Add period if the quote doesn't end with punctuation
	if !strings.HasSuffix(quote, ".") && !strings.HasSuffix(quote, "!") && !strings.HasSuffix(quote, "?") {
		quote += "."
	}

	DebugLog("TextSource: Parsed quote - Content: %s, Author: %s", quote, author)
	return fmt.Sprintf("%s - %s", quote, author), nil
}

func (s *ZenQuotesSource) FormatText(text string) string {
	// Remove any special characters and normalize spaces
	text = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == ' ' {
			return r
		}
		return ' '
	}, text)

	// Normalize spaces
	text = strings.Join(strings.Fields(text), " ")

	// Format based on game mode
	switch CurrentSettings.GameMode {
	case GameModeSimple:
		// For simple mode, convert to lowercase and remove punctuation
		text = strings.ToLower(text)
		words := strings.Fields(text)
		if len(words) > 50 {
			words = words[:50]
		}
		text = strings.Join(words, " ")
	case GameModeNormal:
		// For normal mode, keep the text as is but ensure it's not too long
		words := strings.Fields(text)
		if len(words) > 100 {
			words = words[:100]
		}
		text = strings.Join(words, " ")
	}

	return text
}

type BibleSource struct {
	URL string
}

func NewBibleSource() *BibleSource {
	return &BibleSource{
		URL: "https://bible-api.com/john+3:16",
	}
}

func (s *BibleSource) FetchText() (string, error) {
	DebugLog("TextSource: Fetching bible verse from %s", s.URL)
	resp, err := http.Get(s.URL)
	if err != nil {
		DebugLog("TextSource: Failed to fetch bible verse: %v", err)
		return "", fmt.Errorf("failed to fetch bible verse: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		DebugLog("TextSource: Failed to read response: %v", err)
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	DebugLog("TextSource: Raw API response: %s", string(body))

	var result struct {
		Text      string `json:"text"`
		Reference string `json:"reference"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		DebugLog("TextSource: Failed to parse bible verse: %v", err)
		return "", fmt.Errorf("failed to parse bible verse: %w", err)
	}

	// Clean up the verse text (remove newlines and extra spaces)
	verse := strings.TrimSpace(result.Text)
	verse = strings.Join(strings.Fields(verse), " ")

	DebugLog("TextSource: Parsed verse - Text: %s", verse)
	return verse, nil
}

func (s *BibleSource) FormatText(text string) string {
	// Remove any special characters and normalize spaces
	text = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == ' ' {
			return r
		}
		return ' '
	}, text)

	// Normalize spaces
	text = strings.Join(strings.Fields(text), " ")

	// Format based on game mode
	switch CurrentSettings.GameMode {
	case GameModeSimple:
		// For simple mode, convert to lowercase and remove punctuation
		text = strings.ToLower(text)
		words := strings.Fields(text)
		if len(words) > 50 {
			words = words[:50]
		}
		text = strings.Join(words, " ")
	case GameModeNormal:
		// For normal mode, keep the text as is but ensure it's not too long
		words := strings.Fields(text)
		if len(words) > 100 {
			words = words[:100]
		}
		text = strings.Join(words, " ")
	}

	return text
}

// GetRandomText fetches a random text passage based on the current settings
func GetRandomText() string {
	var source TextSource
	var err error
	var text string

	// Try both sources until we get a successful response
	for i := 0; i < 2; i++ {
		switch i {
		case 0:
			source = NewZenQuotesSource()
			DebugLog("TextSource: Trying ZenQuotes API")
		case 1:
			source = NewBibleSource()
			DebugLog("TextSource: Trying Bible API")
		}

		text, err = source.FetchText()
		if err == nil {
			DebugLog("TextSource: Successfully fetched text: %s", text)
			return source.FormatText(text)
		}
		DebugLog("TextSource: Failed to fetch from source %d: %v", i, err)
	}

	// If all sources fail, return a default text
	DebugLog("TextSource: All sources failed, using default text")
	return "The quick brown fox jumps over the lazy dog."
}
