package ui

import (
	"fmt"
	"strings"

	"goparselogs/internal/models"

	"github.com/charmbracelet/lipgloss"
)

// renderCountdownDisplayView renders the countdown progress modal
func renderCountdownDisplayView(m models.Model) string {
	var view strings.Builder

	view.WriteString("Countdown in progress:\n\n")

	// Visual representation of countdown
	bar := lipgloss.NewStyle().
		Width(30).
		Background(lipgloss.Color("63")).
		Foreground(lipgloss.Color("15")).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("%d seconds remaining", m.CountdownValue))

	view.WriteString(bar + "\n\n")
	view.WriteString("Press ESC to cancel\n")

	// Center the display with a border
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder(), true).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)

	return lipgloss.Place(
		m.TermWidth,
		m.TermHeight,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(view.String()),
	)
}
