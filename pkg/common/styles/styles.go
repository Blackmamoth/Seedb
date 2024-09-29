package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var SuccessStyle lipgloss.Style
var ErrorStyle lipgloss.Style

func init() {
	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Bold(true).
		Italic(true).
		Padding(0, 1).
		Margin(0, 1)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true).
		Underline(true).
		Italic(true).
		Padding(0, 1).
		Margin(0, 1)
}
