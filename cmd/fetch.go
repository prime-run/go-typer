package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// formatText cleans and formats text for typing practice
func formatText(text string) string {
	// Remove extra whitespace
	text = strings.TrimSpace(text)

	// Replace multiple spaces with single space
	text = strings.Join(strings.Fields(text), " ")

	// Remove any non-printable characters
	text = strings.Map(func(r rune) rune {
		if r < 32 || r > 126 {
			return -1
		}
		return r
	}, text)

	return text
}

// formatForGameMode formats text for different game modes
func formatForGameMode(text string, mode string) string {
	text = formatText(text)

	switch mode {
	case "words":
		// Split into words and rejoin with newlines
		words := strings.Fields(text)
		return strings.Join(words, "\n")
	case "sentences":
		// Split on sentence endings and clean up
		text = strings.ReplaceAll(text, ".", ".\n")
		text = strings.ReplaceAll(text, "!", "!\n")
		text = strings.ReplaceAll(text, "?", "?\n")
		// Clean up any double newlines
		lines := strings.Split(text, "\n")
		var cleanLines []string
		for _, line := range lines {
			if clean := strings.TrimSpace(line); clean != "" {
				cleanLines = append(cleanLines, clean)
			}
		}
		return strings.Join(cleanLines, "\n")
	default:
		return text
	}
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Test text fetching from APIs",
	Run: func(cmd *cobra.Command, args []string) {
		modes := []string{"default", "words", "sentences"}

		fmt.Println("Trying ZenQuotes API...")
		zenQuotes := &TextSource{
			URL: "https://zenquotes.io/api/random",
			Parser: func(body []byte) (string, error) {
				var result []struct {
					Quote  string `json:"q"`
					Author string `json:"a"`
				}
				if err := json.Unmarshal(body, &result); err != nil {
					return "", fmt.Errorf("failed to parse JSON: %w", err)
				}
				if len(result) == 0 {
					return "", fmt.Errorf("no quotes found in response")
				}

				quote := formatText(result[0].Quote)
				author := formatText(result[0].Author)

				// Format: "Quote." - Author
				formattedQuote := quote
				if !strings.HasSuffix(quote, ".") && !strings.HasSuffix(quote, "!") && !strings.HasSuffix(quote, "?") {
					formattedQuote += "."
				}
				return fmt.Sprintf("%s - %s", formattedQuote, author), nil
			},
		}
		if text, err := zenQuotes.FetchText(); err != nil {
			fmt.Printf("ZenQuotes API failed: %v\n", err)
		} else {
			fmt.Printf("\nZenQuotes API success:\n")
			for _, mode := range modes {
				formatted := formatForGameMode(text, mode)
				fmt.Printf("\nMode: %s\n", mode)
				fmt.Printf("Text:\n%s\n", formatted)
				fmt.Printf("Character count: %d\n", len(formatted))
				fmt.Printf("Line count: %d\n", len(strings.Split(formatted, "\n")))
			}
		}

		fmt.Println("\nTrying Bible API...")
		bible := &TextSource{
			URL: "https://bible-api.com/john+3:16",
			Parser: func(body []byte) (string, error) {
				var result struct {
					Text      string `json:"text"`
					Reference string `json:"reference"`
				}
				if err := json.Unmarshal(body, &result); err != nil {
					return "", fmt.Errorf("failed to parse JSON: %w", err)
				}

				verse := formatText(result.Text)
				// Don't include the reference for typing practice
				return verse, nil
			},
		}
		if text, err := bible.FetchText(); err != nil {
			fmt.Printf("Bible API failed: %v\n", err)
		} else {
			fmt.Printf("\nBible API success:\n")
			for _, mode := range modes {
				formatted := formatForGameMode(text, mode)
				fmt.Printf("\nMode: %s\n", mode)
				fmt.Printf("Text:\n%s\n", formatted)
				fmt.Printf("Character count: %d\n", len(formatted))
				fmt.Printf("Line count: %d\n", len(strings.Split(formatted, "\n")))
			}
		}
	},
}

// TextSource represents a source for fetching text
type TextSource struct {
	URL    string
	Parser func([]byte) (string, error)
}

// FetchText fetches text from the API
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

func init() {
	rootCmd.AddCommand(fetchCmd)
}
