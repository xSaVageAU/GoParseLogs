package ui

import (
	"fmt"
	"strings"

	"goparselogs/internal/models"
)

// renderMacroListView renders the list of available macros for the right pane.
func renderMacroListView(m models.Model, paneWidth int) string {
	var view strings.Builder

	view.WriteString("Available Macros (UP/DOWN, ENTER):\n\n")
	if len(m.MacroChoices) == 0 {
		view.WriteString(m.SubtleStyle.Render("  No macros defined yet."))
		return m.RightPaneStyle.Copy().Width(paneWidth).Render(view.String())
	}

	for i, choice := range m.MacroChoices {
		cursor := "  "
		line := choice

		// Truncate long macro names
		// Ensure RightPaneStyle is accessible or passed if needed for padding calculation
		// Assuming m.RightPaneStyle.GetHorizontalPadding() is available or we use a const
		availableWidth := paneWidth - m.RightPaneStyle.GetHorizontalPadding() - len(cursor)
		if availableWidth < 5 { // Minimum width for "..." + a couple of chars
			availableWidth = 5
		}

		if len(line) > availableWidth {
			line = line[:availableWidth-3] + "..."
		}

		if m.FocusedPane == models.MacroListPane && m.MacroCursor == i {
			cursor = "> "
			line = m.HighlightStyle.Render(line)
		}
		view.WriteString(fmt.Sprintf("%s%s\n", cursor, line))
	}
	return m.RightPaneStyle.Copy().Width(paneWidth).Render(view.String())
}
