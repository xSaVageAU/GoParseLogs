package ui

import (
	"fmt"
	"strings"

	"goparselogs/internal/models"
)

// renderMenuView renders the menu view with log files and filter input
func renderMenuView(m models.Model) string {
	var view strings.Builder

	// Build menu content
	view.WriteString("Log Files (UP/DOWN, ENTER):\n\n")
	for i, choice := range m.MenuChoices {
		cursor := "  "
		line := choice

		// Truncate long filenames
		availableWidthForFilename := m.TermWidth - m.LeftPaneStyle.GetHorizontalPadding() - len(cursor)
		if availableWidthForFilename < 5 {
			availableWidthForFilename = 5
		}

		if len(line) > availableWidthForFilename {
			line = line[:availableWidthForFilename-3] + "..."
		}

		if m.FocusedPane == models.LogFilePane && m.MenuCursor == i {
			cursor = "> "
			line = m.HighlightStyle.Render(line)
		}
		view.WriteString(fmt.Sprintf("%s%s\n", cursor, line))
	}

	// Filters section
	view.WriteString("\n\nActive Filters:\n")
	if len(m.Filters) == 0 {
		view.WriteString(m.SubtleStyle.Render("  None\n"))
	} else {
		for _, f := range m.Filters {
			view.WriteString(fmt.Sprintf("  - %s\n", f))
		}
	}

	// Filter input section
	filterPrompt := "\nAdd Filter (Type & ENTER):\n"
	if m.CoreProtectMode {
		view.WriteString(m.SubtleStyle.Render("\nFilters disabled in CoreProtect mode.\n"))
	} else {
		if m.TermWidth < 45 {
			filterPrompt = "\nFilter:\n"
		}
		view.WriteString(filterPrompt)
		currentInputStyle := m.InputStyle
		filterText := m.FilterInput
		if m.FocusedPane == models.FilterPane {
			currentInputStyle = m.FocusedInputStyle
			filterText += "▌"
		}

		inputRenderWidth := m.TermWidth - m.LeftPaneStyle.GetHorizontalPadding() - currentInputStyle.GetHorizontalFrameSize() - 2
		if inputRenderWidth < 5 {
			inputRenderWidth = 5
		}
		view.WriteString(currentInputStyle.Width(inputRenderWidth).Render(filterText))
	}

	// Help text
	helpText := buildHelpText(m)
	view.WriteString(m.SubtleStyle.Render(helpText))

	// Error display
	if m.Err != nil {
		view.WriteString("\n\n" + m.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.Err)))
	}

	// Save message display
	if m.SaveMessage != "" {
		styleToUse := m.SubtleStyle
		if strings.HasPrefix(strings.ToLower(m.SaveMessage), "error") {
			styleToUse = m.ErrorStyle
		} else if strings.HasPrefix(strings.ToLower(m.SaveMessage), "logs saved") {
			styleToUse = m.SuccessStyle
		}
		view.WriteString("\n\n" + styleToUse.Render(m.SaveMessage))
	}

	return m.LeftPaneStyle.Copy().Width(m.TermWidth - m.LeftPaneStyle.GetHorizontalFrameSize()).Render(view.String())
}
