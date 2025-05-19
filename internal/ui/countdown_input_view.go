package ui

import (
	"strings"

	"goparselogs/internal/models"

	"github.com/charmbracelet/lipgloss"
)

// renderCountdownInputView renders the countdown input modal
func renderCountdownInputView(m models.Model) string {
	var view strings.Builder

	view.WriteString("Enter countdown time in seconds (ENTER to start, ESC to cancel):\n\n")

	// Use a style for the input field without its own border to avoid conflict with modal border
	inputRenderStyle := m.FocusedInputStyle.Copy().Border(lipgloss.Border{})
	view.WriteString(inputRenderStyle.Width(m.TermWidth / 4).Render(m.CountdownInput + "▌"))
	view.WriteString("\n\n")

	if m.CountdownMessage != "" {
		styleToUse := m.SubtleStyle
		if strings.HasPrefix(strings.ToLower(m.CountdownMessage), "invalid") {
			styleToUse = m.ErrorStyle
		}
		view.WriteString(styleToUse.Render(m.CountdownMessage) + "\n")
	}

	// Center the input box with a border
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
