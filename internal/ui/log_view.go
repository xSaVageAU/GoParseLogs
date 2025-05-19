package ui

import (
	"fmt"
	"strings"

	"goparselogs/internal/models"

	"github.com/charmbracelet/lipgloss"
)

// renderLogView renders the split view showing both menu and logs
func renderLogView(m models.Model) string {
	// Build left pane (menu) content
	var leftPane strings.Builder

	leftPane.WriteString("Log Files (UP/DOWN, ENTER):\n\n")
	for i, choice := range m.MenuChoices {
		cursor := "  "
		line := choice

		availableWidthForFilename := m.LeftPaneWidth - m.LeftPaneStyle.GetHorizontalPadding() - len(cursor)
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
		leftPane.WriteString(fmt.Sprintf("%s%s\n", cursor, line))
	}

	// Filters section
	leftPane.WriteString("\n\nActive Filters:\n")
	if len(m.Filters) == 0 {
		leftPane.WriteString(m.SubtleStyle.Render("  None\n"))
	} else {
		for _, f := range m.Filters {
			leftPane.WriteString(fmt.Sprintf("  - %s\n", f))
		}
	}

	// Filter input
	filterPrompt := "\nAdd Filter (Type & ENTER):\n"
	if m.CoreProtectMode {
		leftPane.WriteString(m.SubtleStyle.Render("\nFilters disabled in CoreProtect mode.\n"))
	} else {
		if m.LeftPaneWidth < 45 {
			filterPrompt = "\nFilter:\n"
		}
		leftPane.WriteString(filterPrompt)
		currentInputStyle := m.InputStyle
		filterText := m.FilterInput
		if m.FocusedPane == models.FilterPane {
			currentInputStyle = m.FocusedInputStyle
			filterText += "▌"
		}

		inputRenderWidth := m.LeftPaneWidth - m.LeftPaneStyle.GetHorizontalPadding() - currentInputStyle.GetHorizontalFrameSize() - 2
		if inputRenderWidth < 5 {
			inputRenderWidth = 5
		}
		leftPane.WriteString(currentInputStyle.Width(inputRenderWidth).Render(filterText))
	}

	// Help text
	leftPane.WriteString(m.SubtleStyle.Render(buildHelpText(m)))

	// Error and save message display
	if m.Err != nil {
		leftPane.WriteString("\n\n" + m.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.Err)))
	}
	if m.SaveMessage != "" {
		styleToUse := m.SubtleStyle
		if strings.HasPrefix(strings.ToLower(m.SaveMessage), "error") {
			styleToUse = m.ErrorStyle
		} else if strings.HasPrefix(strings.ToLower(m.SaveMessage), "logs saved") {
			styleToUse = m.SuccessStyle
		}
		leftPane.WriteString("\n\n" + styleToUse.Render(m.SaveMessage))
	}

	styledLeftPane := m.LeftPaneStyle.Width(m.LeftPaneWidth).Render(leftPane.String())

	// Build right pane (logs) content
	var rightPane strings.Builder
	rightPaneWidth := 0
	if m.TermWidth > m.LeftPaneWidth+m.LeftPaneStyle.GetHorizontalBorderSize() {
		rightPaneWidth = m.TermWidth - m.LeftPaneWidth - m.LeftPaneStyle.GetHorizontalBorderSize()
	}

	if rightPaneWidth <= 10 {
		rightPane.WriteString(m.ErrorStyle.Render("Terminal too narrow."))
	} else if m.State == models.MacroListView {
		// Call the dedicated function to render the macro list
		// Note: We pass 'rightPane' builder directly to styledRightPane later.
		// So, renderMacroListView should return the string to be set for styledRightPane.
		// For now, we'll let renderMacroListView handle its own styling and width.
		// This means the styledRightPane will effectively be the result of renderMacroListView.
		// This is a slight departure but simplifies the call here.
		// The alternative is renderMacroListView writes to 'rightPane' builder,
		// but then it can't apply its own m.RightPaneStyle.Copy().Width().Render().
		// Let's assume renderMacroListView returns a fully styled string.
		// We will assign its result to styledRightPane directly later.
		// For now, this 'else if' block for MacroListView will be handled differently below.
	} else if m.State == models.MenuView {
		// Custom message when no file is selected
		rightPane.WriteString("Select a log file from the left panel to view its contents.\n\n")
		rightPane.WriteString(m.SubtleStyle.Render("Use UP/DOWN or J/K to navigate\n"))
		rightPane.WriteString(m.SubtleStyle.Render("Press ENTER to view a file\n"))
		if !m.CoreProtectMode {
			rightPane.WriteString(m.SubtleStyle.Render("Press TAB to focus on filters\n"))
		}
	} else if len(m.LogEntries) == 0 && m.Err == nil && !m.CoreProtectMode {
		if len(m.Filters) > 0 {
			rightPane.WriteString(fmt.Sprintf("No log entries matching filters: %s\n", strings.Join(m.Filters, ", ")))
		} else {
			rightPane.WriteString("Loading or parsing log file...")
		}
	} else if len(m.CoreProtectLogEntries) == 0 && m.Err == nil && m.CoreProtectMode {
		rightPane.WriteString("Loading or parsing CoreProtect log file...")
	} else if m.Err != nil {
		rightPane.WriteString(m.ErrorStyle.Render("Error loading logs. See left pane."))
	} else {
		headerFooterAndPaddingHeight := m.RightPaneStyle.GetVerticalPadding() + 2 + 1 + 1 + 1 + 1
		availableHeightForLogs := m.TermHeight - headerFooterAndPaddingHeight
		if availableHeightForLogs < 1 {
			availableHeightForLogs = 1
		}
		numEntriesToShow := availableHeightForLogs

		var currentEntriesCount int
		if m.CoreProtectMode {
			rightPane.WriteString("CoreProtect Log Entries (Sorted by Hours Ago):\n\n")
			currentEntriesCount = len(m.CoreProtectLogEntries)
		} else {
			rightPane.WriteString("Parsed Log Entries")
			if len(m.Filters) > 0 {
				rightPane.WriteString(fmt.Sprintf(" (Filters: %s)", m.HighlightStyle.Render(strings.Join(m.Filters, ", "))))
			}
			rightPane.WriteString(":\n\n")
			currentEntriesCount = len(m.LogEntries)
		}

		if currentEntriesCount == 0 {
			if m.CoreProtectMode {
				rightPane.WriteString("No CoreProtect entries found or parsed.")
			} else if len(m.Filters) > 0 {
				rightPane.WriteString(fmt.Sprintf("No log entries matching filters: %s", strings.Join(m.Filters, ", ")))
			} else {
				rightPane.WriteString("No log entries.")
			}
		} else {
			start := m.LogCursor - (numEntriesToShow / 2)
			if start < 0 {
				start = 0
			}
			end := start + numEntriesToShow
			if end > currentEntriesCount {
				end = currentEntriesCount
				start = end - numEntriesToShow
				if start < 0 {
					start = 0
				}
			}
			if end == start && currentEntriesCount > 0 {
				end = start + 1
			}

			for i := start; i < end; i++ {
				var line string
				if m.CoreProtectMode {
					entry := m.CoreProtectLogEntries[i]
					line = fmt.Sprintf("%.2f/h ago - %s: %s", entry.HoursAgo, entry.Username, entry.Message)
				} else {
					entry := m.LogEntries[i]
					line = fmt.Sprintf("[%s] [%s/%s]: %s", entry.Timestamp, entry.Thread, entry.Level, entry.Message)
				}

				maxLineTextWidth := rightPaneWidth - m.RightPaneStyle.GetHorizontalPadding() - 2
				if maxLineTextWidth < 5 {
					maxLineTextWidth = 5
				}

				if len(line) > maxLineTextWidth {
					line = line[:maxLineTextWidth-3] + "..."
				}

				var styledLine string
				if i == m.LogCursor {
					styledLine = m.HighlightStyle.Render(fmt.Sprintf("> %s", line))
				} else {
					styledLine = fmt.Sprintf("  %s", line)
				}
				rightPane.WriteString(styledLine + "\n")
			}

			if currentEntriesCount > 0 {
				rightPane.WriteString(fmt.Sprintf("\nViewing %d-%d of %d\n", start+1, end, currentEntriesCount))
			} else {
				rightPane.WriteString("\nNo entries to display.\n")
			}
		}
	}

	var styledRightPane string
	if m.State == models.MacroListView {
		styledRightPane = renderMacroListView(m, rightPaneWidth)
	} else {
		styledRightPane = m.RightPaneStyle.Width(rightPaneWidth).Render(rightPane.String())
	}

	// Combine panes horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, styledLeftPane, styledRightPane)
}
