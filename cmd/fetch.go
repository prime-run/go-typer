// NOTE: maybe cache using a free edge run-time and store in an S3 bucket?
package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/prime-run/go-typer/types"
	"github.com/prime-run/go-typer/utils"
	"github.com/spf13/cobra"
)

const (
	ModeDefault   types.Mode = "default"
	ModeWords     types.Mode = "words"
	ModeSentences types.Mode = "sentences"

	zenQuotesAPIURL = "https://zenquotes.io/api/random"
	bibleAPIURL     = "https://bible-api.com/john+3:16"
)

func init() {
	rootCmd.AddCommand(fetchCmd)
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Test text fetching from APIs",
	Run: func(cmd *cobra.Command, args []string) {
		modes := []types.Mode{ModeDefault, ModeWords, ModeSentences}

		fmt.Println("Trying ZenQuotes API...")
		zenQuotes := &types.TextSource{
			URL: zenQuotesAPIURL,
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

				quote := utils.FormatText(result[0].Quote)
				author := utils.FormatText(result[0].Author)

				formattedQuote := quote
				if utils.HasPonctuationSuffix(quote) {
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
				utils.PrintTextStats(mode, formatted)
			}
		}

		fmt.Println("\nTrying Bible API...")
		bible := &types.TextSource{
			URL: bibleAPIURL,
			Parser: func(body []byte) (string, error) {
				var result struct {
					Text string `json:"text"`
				}
				if err := json.Unmarshal(body, &result); err != nil {
					return "", fmt.Errorf("failed to parse JSON: %w", err)
				}

				verse := utils.FormatText(result.Text)
				// !WARN: don't include the reference for typing practice
				return verse, nil
			},
		}
		if text, err := bible.FetchText(); err != nil {
			fmt.Printf("Bible API failed: %v\n", err)
		} else {
			fmt.Printf("\nBible API success:\n")
			for _, mode := range modes {
				formatted := formatForGameMode(text, mode)
				utils.PrintTextStats(mode, formatted)
			}
		}
	},
}

func formatForGameMode(text string, mode types.Mode) string {
	text = utils.FormatText(text)

	switch mode {
	case ModeWords:
		return formatForWords(text)
	case ModeSentences:
		return formatForSentences(text)
	default:
		return text
	}
}

func formatForWords(text string) string {
	words := strings.Fields(text)
	return strings.Join(words, "\n")
}

func formatForSentences(text string) string {
	text = strings.ReplaceAll(text, ".", ".\n")
	text = strings.ReplaceAll(text, "!", "!\n")
	text = strings.ReplaceAll(text, "?", "?\n")

	var cleanLines []string
	for line := range strings.Lines(text) {
		if clean := strings.TrimSpace(line); clean != "" {
			cleanLines = append(cleanLines, clean)
		}
	}
	return strings.Join(cleanLines, "\n")
}
