package ui

import (
	"strings"

	"goparselogs/internal/models"
)

// buildHelpText creates the appropriate help text based on the current state
func buildHelpText(m models.Model) string {
	var helpText string
	var helpParts []string

	baseHelp := []string{"TAB: Focus", "Q/^C: Quit"}

	switch m.State {
	case models.LogView:
		specificHelp := []string{"E: Save", "ESC: Menu"}
		if m.LeftPaneWidth < 45 { // Threshold for single line help
			helpParts = append(baseHelp, specificHelp...)
			helpText = "\n" + strings.Join(helpParts, " | ")
		} else {
			helpText = "\n" + strings.Join(baseHelp, " | ") + "\n" + strings.Join(specificHelp, " | ")
		}

	case models.MenuView:
		specificHelp := []string{"ESC: Unfocus"}
		if m.LeftPaneWidth < 40 {
			helpParts = append(baseHelp, specificHelp...)
			helpText = "\n" + strings.Join(helpParts, " | ")
		} else {
			helpText = "\n" + strings.Join(baseHelp, " | ") + "\n" + strings.Join(specificHelp, " | ")
		}

	case models.SaveInputView:
		helpText = "\nEnter filename. ENTER: Save, ESC: Cancel."
	}

	return helpText
}
