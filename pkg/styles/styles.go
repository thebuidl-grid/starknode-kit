package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Primary style – used for headlines or main buttons
	Primary = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)

	// Secondary style – used for secondary buttons or muted text
	Secondary = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c757d")). // gray
			Padding(0, 1)

	// Danger style – for errors
	Danger = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#dc3545")). // red
		Padding(0, 0)

	// Warning style – for caution/warnings
	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffc107")). // yellow
		Padding(0, 1)

	// Info style – for informational banners
	Info = lipgloss.NewStyle().
		Background(lipgloss.Color("#17a2b8")). // cyan
		Padding(0, 1)

	// Muted text – low-priority, descriptive
	Muted = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c757d"))

	// Title style – bold, prominent, no blue foreground
	Title = lipgloss.NewStyle().
		Bold(true).
		Underline(true)

	// Subtitle style
	Subtitle = lipgloss.NewStyle().
			Bold(true)

	// BorderBox style – e.g. boxed alerts
	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder())
)

// Checkbox renders a styled checkbox for menu items
func Checkbox(label string, checked bool) string {
	if checked {
		return Danger.Render("[✔]") + " " + label
	}
	return "[ ] " + label
}
