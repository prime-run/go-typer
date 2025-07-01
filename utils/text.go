package utils

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/prime-run/go-typer/types"
)

func FormatText(text string) string {
	text = sanitizeText(text)
	return text
}

func sanitizeText(text string) string {
	text = strings.Map(func(r rune) rune {
		if !unicode.IsPrint(r) {
			return -1
		}
		return r
	}, strings.TrimSpace(text))

	text = strings.Join(strings.Fields(text), " ")

	return text
}

func PrintTextStats(mode types.Mode, formatedTxt string) {
	fmt.Printf("\nMode: %s\n", mode)
	fmt.Printf("Text:\n%s\n", formatedTxt)
	fmt.Printf("Character count: %d\n", len(formatedTxt))
	fmt.Printf("Line count: %d\n", strings.Count(formatedTxt, "\n")+1)
}

func HasPonctuationSuffix(text string) bool {
	if len(text) == 0 {
		return false
	}
	return strings.HasSuffix(text, ".") || strings.HasSuffix(text, "!") || strings.HasSuffix(text, "?")
}
