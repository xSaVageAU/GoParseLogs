package ui

import (
	"fmt"
	"strings"

	"goparselogs/internal/models"
)

// View renders the entire UI based on the current state
func View(m models.Model) string {
	if m.TermWidth == 0 { // Wait for initial WindowSizeMsg
		return "Initializing..."
	}

	// Update CoreProtect toggle text in menuChoices
	for i, choice := range m.MenuChoices {
		if strings.HasPrefix(choice, CoreProtectToggleBaseText) {
			statusText := "OFF"
			if m.CoreProtectMode {
				statusText = "ON"
			}
			m.MenuChoices[i] = fmt.Sprintf("%s (%s)", CoreProtectToggleBaseText, statusText)
			break
		}
	}

	var finalView strings.Builder

	switch m.State {
	case models.SaveInputView:
		// For save dialog, we use a modal overlay
		finalView.WriteString(renderLogView(m)) // Render the background
		finalView.WriteString("\n")
		finalView.WriteString(renderSaveInputView(m)) // Overlay the save dialog
	default:
		// For all other states, use the split view
		finalView.WriteString(renderLogView(m))
	}

	return finalView.String()
}
