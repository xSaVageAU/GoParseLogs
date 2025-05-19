package ui

import (
	"strings"

	"goparselogs/internal/models"

	"github.com/charmbracelet/lipgloss"
)

// renderSaveInputView renders the save file input modal
func renderSaveInputView(m models.Model) string {
	var saveView strings.Builder

	saveView.WriteString("Enter filename to save logs (ENTER to save, ESC to cancel):\n\n")

	// Use a style for the save input field without its own border to avoid conflict with modal border
	saveInputRenderStyle := m.FocusedInputStyle.Copy().Border(lipgloss.Border{})
	saveView.WriteString(saveInputRenderStyle.Width(m.TermWidth / 2).Render(m.SaveFilenameInput + "▌"))
	saveView.WriteString("\n\n")

	if m.SaveMessage != "" {
		styleToUse := m.SubtleStyle
		if strings.HasPrefix(strings.ToLower(m.SaveMessage), "error") {
			styleToUse = m.ErrorStyle
		}
		saveView.WriteString(styleToUse.Render(m.SaveMessage) + "\n")
	}

	// Center the save input box with a border
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder(), true).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)

	return lipgloss.Place(
		m.TermWidth,
		m.TermHeight,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(saveView.String()),
	)
}
