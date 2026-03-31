package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorBackground = lipgloss.Color("#090040")
	ColorSecondary  = lipgloss.Color("#471396")
	ColorAccent     = lipgloss.Color("#B13BFF")
	ColorHighlight  = lipgloss.Color("#FFCC00")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Background(ColorBackground).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Background(ColorBackground).
			Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Bold(true).
			MarginBottom(1)

	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorHighlight)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Background(ColorBackground).
			Padding(1, 2)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)
)
