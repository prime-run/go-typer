package ui

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

var GradientColors = []string{
	"#00ADD8", // go blue
	"#15B5DB",
	"#2ABEDE",
	"#3FC6E1",
	"#54CFE4",
	"#69D7E7",
	"#7EE0EA",
	"#93E8ED",
	"#A8F1F0",
	"#BDF9F3",
	"#D2FFF6",
	"#E7FFF9",
	"#FCFFFC",
	"#E7FFF9",
	"#D2FFF6",
	"#BDF9F3",
	"#A8F1F0",
	"#93E8ED",
	"#7EE0EA",
	"#69D7E7",
	"#54CFE4",
	"#3FC6E1",
	"#2ABEDE",
	"#15B5DB",
}

func GetGradientIndex(tickTime time.Time) int {
	return int(tickTime.UnixNano()/int64(30*time.Millisecond)) % len(GradientColors)
}

func RenderGradientText(text string, tickTime time.Time) string {
	var result strings.Builder
	colorIndex := GetGradientIndex(tickTime)

	result.Grow(len(text) * 3)

	for _, char := range text {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(GradientColors[colorIndex]))
		result.WriteString(style.Render(string(char)))
		colorIndex = (colorIndex + 1) % len(GradientColors)
	}
	return result.String()
}

func RenderGradientOverlay(text string, baseStyle lipgloss.Style, tickTime time.Time) string {
	var result strings.Builder
	colorIndex := GetGradientIndex(tickTime)

	result.Grow(len(text) * 3)

	for _, char := range text {
		gradientColor := lipgloss.Color(GradientColors[colorIndex])

		combinedStyle := baseStyle.Copy().
			Foreground(gradientColor).
			Bold(true)

		result.WriteString(combinedStyle.Render(string(char)))
		colorIndex = (colorIndex + 1) % len(GradientColors)
	}

	return result.String()
}
