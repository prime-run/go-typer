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

var HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

// var TextToTypeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#333333")).Padding(1).Width(MaxWidth)
var TextToTypeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7F9ABE")).Padding(1).Width(MaxWidth)
var InputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
var ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
var CenterStyle = lipgloss.NewStyle().Align(lipgloss.Center)

// Padding style for consistent spacing
var PadStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

// Timer display style
var TimerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFDB58")).
	Bold(true).
	Padding(0, 1)

// Live typing preview box style
var PreviewStyle = lipgloss.NewStyle().
	Padding(1).
	Margin(8, 0, 0, 0).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#7F9ABE")).
	Width(MaxWidth)

const SampleText = "The quick brown fox jumps over the lazy dog. Programming is the process of creating a set of instructions that tell a computer how to perform a task. Programming can be done using a variety of computer programming languages, such as JavaScript, Python, and C++."
