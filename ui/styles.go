package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	Padding  = 2
	MaxWidth = 80
)

// WARN:switched to true color might comeback to bite later in testing for other termnal emulators!

// TODO:a theming system would be nice

func UpdateStyles() {

	helpStyle := lipgloss.NewStyle().Foreground(GetColor("help_text"))
	HelpStyle = helpStyle.Render

	TextToTypeStyle = lipgloss.NewStyle().Foreground(GetColor("text_preview")).Padding(1).Width(MaxWidth)
	InputStyle = lipgloss.NewStyle().Foreground(GetColor("text_correct"))
	ErrorStyle = lipgloss.NewStyle().Foreground(GetColor("text_error"))
	PartialErrorStyle = lipgloss.NewStyle().Foreground(GetColor("text_partial_error"))
	DimStyle = lipgloss.NewStyle().Foreground(GetColor("text_dim"))

	CenterStyle = lipgloss.NewStyle().Align(lipgloss.Center)
	PadStyle = lipgloss.NewStyle().Foreground(GetColor("padding"))

	TimerStyle = lipgloss.NewStyle().
		Foreground(GetColor("timer")).
		Bold(true).
		Padding(0, 1)

	PreviewStyle = lipgloss.NewStyle().
		Padding(1).
		Margin(8, 0, 0, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GetColor("border")).
		Width(MaxWidth)

	TextContainerStyle = lipgloss.NewStyle().
		Padding(1).
		Width(MaxWidth)

	BlockCursorStyle = lipgloss.NewStyle().
		Foreground(GetColor("cursor_fg")).
		Background(GetColor("cursor_bg"))

	UnderlineCursorStyle = lipgloss.NewStyle().
		Foreground(GetColor("cursor_underline")).
		Underline(true)
}

var HelpStyle func(...string) string
var TextToTypeStyle lipgloss.Style
var InputStyle lipgloss.Style
var ErrorStyle lipgloss.Style
var PartialErrorStyle lipgloss.Style
var CenterStyle lipgloss.Style
var PadStyle lipgloss.Style
var TimerStyle lipgloss.Style
var PreviewStyle lipgloss.Style
var DimStyle lipgloss.Style
var TextContainerStyle lipgloss.Style
var BlockCursorStyle lipgloss.Style
var UnderlineCursorStyle lipgloss.Style

const (
	SampleTextNormal = "The quick brown fox jumps over the lazy dog. Programming is the process of creating a set of instructions that tell a computer how to perform a task. Programming can be done using a variety of computer programming languages, such as JavaScript, Python, and C++."

	SampleTextNormalWithNumbers = "The quick brown fox jumps over the 5 lazy dogs. In 2023, programming is the process of creating a set of instructions that tell a computer how to perform a task. Programming can be done using a variety of computer programming languages, such as JavaScript, Python, and C++, with over 300 languages in existence."

	SampleTextSimple = "the quick brown fox jumps over the lazy dog programming is the process of creating a set of instructions that tell a computer how to perform a task programming can be done using a variety of computer programming languages such as javascript python and c plus plus"

	SampleTextSimpleWithNumbers = "the quick brown fox jumps over 5 lazy dogs in 2023 programming is the process of creating a set of instructions that tell a computer how to perform a task programming can be done using a variety of computer programming languages such as javascript python and c plus plus with over 300 languages in existence"
)

func init() {

	InitTheme()

	UpdateStyles()
}

func GetSampleText() string {
	if CurrentSettings.GameMode == GameModeSimple {
		if CurrentSettings.UseNumbers {
			return SampleTextSimpleWithNumbers
		}
		return SampleTextSimple
	} else {
		if CurrentSettings.UseNumbers {
			return SampleTextNormalWithNumbers
		}
		return SampleTextNormal
	}
}
