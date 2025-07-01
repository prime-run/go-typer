package utils

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/prime-run/go-typer/types"
)

// formatText is used to format the text for display.
func FormatText(text string) string {
	text = sanitizeText(text)
	return text
}

// sanitizeText is used to clean up the text before displaying it.
// It removes non-printable characters and replaces multiple spaces with a single space.
func sanitizeText(text string) string {
	// Remove non-printable characters
	text = strings.Map(func(r rune) rune {
		if !unicode.IsPrint(r) {
			return -1
		}
		return r
	}, strings.TrimSpace(text))

	// Replace multiple spaces with a single space
	text = strings.Join(strings.Fields(text), " ")

	return text
}

// PrintTextStats is a utility function to print the formatted text along with its mode and character count.
func PrintTextStats(mode types.Mode, formatedTxt string) {
	fmt.Printf("\nMode: %s\n", mode)
	fmt.Printf("Text:\n%s\n", formatedTxt)
	fmt.Printf("Character count: %d\n", len(formatedTxt))
	fmt.Printf("Line count: %d\n", strings.Count(formatedTxt, "\n")+1)
}

// HasPonctuationSuffix checks if the text ends with a punctuation mark (., !, ?).
func HasPonctuationSuffix(text string) bool {
	if len(text) == 0 {
		return false
	}
	return strings.HasSuffix(text, ".") || strings.HasSuffix(text, "!") || strings.HasSuffix(text, "?")
}
