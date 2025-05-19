package ui

import (
	"fmt"

	"goparselogs/internal/fileops"
	"goparselogs/internal/macros" // Import the new macros package
	"goparselogs/internal/models"
	"goparselogs/pkg/coreprotectparser"
	"goparselogs/pkg/logparser"

	"github.com/charmbracelet/lipgloss"
)

const (
	CoreProtectToggleBaseText = "Toggle CoreProtect Parsing"
	MacrosText                = "Macros"
	ExitText                  = "Exit"
)

// createInitialState creates and returns a new model with default settings and styles
func createInitialState() models.Model {
	// Get list of log files
	logFiles, err := fileops.ScanLogFiles()
	if err != nil {
		// If we can't read the logs directory, start with empty list
		logFiles = []string{}
	}

	// Create menu choices with log files
	menuChoices := make([]string, 0, len(logFiles)+3) // +3 for Macros, toggle, and exit
	menuChoices = append(menuChoices, logFiles...)
	menuChoices = append(menuChoices,
		MacrosText,
		fmt.Sprintf("%s (OFF)", CoreProtectToggleBaseText),
		ExitText,
	)

	highlightStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10"))

	subtleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")). // Dim border for non-focused
		PaddingLeft(1).
		PaddingRight(1).
		Width(50)

	focusedInputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")). // Bright border for focused
		PaddingLeft(1).
		PaddingRight(1).
		Width(50)

	leftPaneStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.NormalBorder(), false, true, false, false). // Border on the right
		BorderForeground(lipgloss.Color("238"))

	rightPaneStyle := lipgloss.NewStyle().
		Padding(1, 2)

	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))    // Red
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green

	return models.Model{
		State:                 models.MenuView,
		FocusedPane:           models.LogFilePane,
		LeftPaneWidth:         60, // Initial default, will be updated by WindowSizeMsg
		MenuChoices:           menuChoices,
		MenuCursor:            0,
		Filters:               []string{},
		CoreProtectMode:       false,
		LogEntries:            []logparser.LogEntry{},
		CoreProtectLogEntries: []coreprotectparser.CoreProtectLogEntry{},
		LogCursor:             0,
		MacroChoices:          macros.GetMacroNames(), // Use the function from the macros package
		MacroCursor:           0,
		MacroParameters:       make(map[string]string), // Initialize empty parameters map
		ParameterInputs:       make(map[string]string), // Initialize empty parameter inputs map
		ParameterCursor:       0,                       // First parameter is selected initially
		ParameterMessage:      "",                      // No parameter message initially
		HighlightStyle:        highlightStyle,
		SubtleStyle:           subtleStyle,
		InputStyle:            inputStyle,
		FocusedInputStyle:     focusedInputStyle,
		LeftPaneStyle:         leftPaneStyle,
		RightPaneStyle:        rightPaneStyle,
		ErrorStyle:            errorStyle,
		SuccessStyle:          successStyle,
		InputActive:           false, // Initially, log file pane is active
	}
}

// Helper functions for min/max
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
