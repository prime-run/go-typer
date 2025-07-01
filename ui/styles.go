package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	Padding  = 2  // Padding for the text
	MaxWidth = 80 // Maximum width for the text

	SampleTextNormal            = "The quick brown fox jumps over the lazy dog. Programming is the process of creating a set of instructions that tell a computer how to perform a task. Programming can be done using a variety of computer programming languages, such as JavaScript, Python, and C++."
	SampleTextNormalWithNumbers = "The quick brown fox jumps over the 5 lazy dogs. In 2023, "
	SampleTextSimple            = "the quick brown fox jumps over the lazy dog programming is the process of creating a set of instructions that tell a computer how to perform a task programming can be done using a variety of computer programming languages such as javascript python and c plus plus"
	SampleTextSimpleWithNumbers = "the quick brown fox jumps over 5 lazy dogs in 2023 programming is the process of creating a set of instructions that tell a computer how to perform a task programming can be done using a variety of computer programming languages such as javascript python and c plus plus with over 300 languages in existence"
)

var (
	HelpStyle                  func(...string) string // Help text style
	HintStyle                  func(...string) string // Hint text style
	SettingsStyle              func(...string) string // Settings text style
	TextToTypeStyle            lipgloss.Style         // Text to type style
	InputStyle                 lipgloss.Style         // Input text style
	ErrorStyle                 lipgloss.Style         // Error text style
	PartialErrorStyle          lipgloss.Style         // Partial error text style
	CenterStyle                lipgloss.Style         // Centered text style
	PadStyle                   lipgloss.Style         // Padding text style
	TimerStyle                 lipgloss.Style         // Timer text style
	PreviewStyle               lipgloss.Style         // Preview text style
	DimStyle                   lipgloss.Style         // Dim text style
	TextContainerStyle         lipgloss.Style         // Text container style
	BlockCursorStyle           lipgloss.Style         // Block cursor style
	UnderlineCursorStyle       lipgloss.Style         // Underline cursor style
	SettingsListStyle          lipgloss.Style         // Settings list style
	SettingsDetailsStyle       lipgloss.Style         // Settings details style
	SettingsTitleStyle         lipgloss.Style         // Settings title style
	SettingsHelpStyle          lipgloss.Style         // Settings help text style
	EndGameTitleStyle          lipgloss.Style         // End game title style
	EndGameStatsBoxStyle       lipgloss.Style         // End game stats box style
	EndGameWpmStyle            lipgloss.Style         // End game WPM style
	EndGameAccuracyStyle       lipgloss.Style         // End game accuracy style
	EndGameWordsStyle          lipgloss.Style         // End game words style
	EndGameCorrectStyle        lipgloss.Style         // End game correct style
	EndGameErrorsStyle         lipgloss.Style         // End game errors style
	EndGameOptionStyle         lipgloss.Style         // End game option style
	EndGameSelectedOptionStyle lipgloss.Style         // End game selected option style
)

// WARN:switched to true color might comeback to bite later in testing for other termnal emulators!

// TODO:a theming system would be nice [x]

func UpdateStyles() {
	helpStyle := lipgloss.NewStyle().Foreground(GetColor("help_text"))
	HelpStyle = helpStyle.Render

	hintStyle := lipgloss.NewStyle().Foreground(GetColor("text_preview")).Italic(true)
	HintStyle = hintStyle.Render

	settingsStyle := lipgloss.NewStyle().Foreground(GetColor("timer")).Bold(true)
	SettingsStyle = settingsStyle.Render

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

	SettingsListStyle = lipgloss.NewStyle().
		Width(MaxWidth/3 - 4).
		MarginLeft(2).
		MarginRight(2)

	SettingsDetailsStyle = lipgloss.NewStyle().
		Width(MaxWidth / 2).
		MarginLeft(2)

	SettingsTitleStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("0")).
		Padding(0, 1)

	SettingsHelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	EndGameTitleStyle = lipgloss.NewStyle().
		Foreground(GetColor("text_correct")).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GetColor("border")).
		Padding(0, 2)

	EndGameStatsBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GetColor("border")).
		Padding(1, 2)

	EndGameWpmStyle = lipgloss.NewStyle().
		Foreground(GetColor("timer")).
		Bold(true).
		Underline(true)

	EndGameAccuracyStyle = lipgloss.NewStyle().
		Foreground(GetColor("text_correct"))

	EndGameWordsStyle = lipgloss.NewStyle().
		Foreground(GetColor("text_preview"))

	EndGameCorrectStyle = lipgloss.NewStyle().
		Foreground(GetColor("text_correct"))

	EndGameErrorsStyle = lipgloss.NewStyle().
		Foreground(GetColor("text_error"))

	EndGameOptionStyle = lipgloss.NewStyle().
		Foreground(GetColor("text_preview"))

	EndGameSelectedOptionStyle = lipgloss.NewStyle().
		Foreground(GetColor("text_correct")).
		Bold(true)
}

// init initializes the styles and themes
func init() {
	InitTheme()
	UpdateStyles()
}

// GetSampleText returns a sample text based on the current settings
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
