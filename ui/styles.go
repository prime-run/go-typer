package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	Padding  = 2
	MaxWidth = 80
)

var HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
var TextToTypeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#333333")).Padding(1).Width(MaxWidth)
var InputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
var ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
var CenterStyle = lipgloss.NewStyle().Align(lipgloss.Center)

const SampleText = "The quick brown fox jumps over the lazy dog. Programming is the process of creating a set of instructions that tell a computer how to perform a task. Programming can be done using a variety of computer programming languages, such as JavaScript, Python, and C++."
