package utils

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func GetThemePath(themeName string) string {
	const YMLSuffix = ".yml"

	themeName = strings.TrimPrefix(themeName, "-")
	if strings.HasSuffix(themeName, YMLSuffix) {
		return themeName
	}

	configDir, err := GetAppConfigDir()
	if err != nil {
		return filepath.Join("colorschemes", themeName+YMLSuffix)
	}

	colorschemesDir := filepath.Join(configDir, "colorschemes")
	if err := os.MkdirAll(colorschemesDir, 0755); err != nil {
		return filepath.Join("colorschemes", themeName+YMLSuffix)
	}

	return filepath.Join(colorschemesDir, themeName+YMLSuffix)
}

func GetDisplayThemeName(themeName string) string {
	if strings.Contains(themeName, "/") || strings.Contains(themeName, "\\") {
		themeName = filepath.Base(themeName)
	}

	themeName = strings.TrimSuffix(themeName, ".yml")

	words := strings.FieldsFunc(themeName, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return strings.Join(words, " ")
}

func IsValidThemeName(name string) bool {
	if strings.Contains(name, ".") && !strings.HasSuffix(name, ".yml") {
		return false
	}

	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		_, err := os.Stat(name)
		return err == nil
	}

	for _, c := range name {
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '-') {
			return false
		}
	}

	return true
}
