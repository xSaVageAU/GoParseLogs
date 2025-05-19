package ui

import (
	"fmt"
	"strings"
	"time"

	"goparselogs/internal/fileops"
	"goparselogs/internal/macros"
	"goparselogs/internal/models"
	"goparselogs/pkg/coreprotectparser"
	"goparselogs/pkg/logparser"

	tea "github.com/charmbracelet/bubbletea"
)

// scanLogsMsg is sent when the logs directory has been rescanned
type scanLogsMsg struct {
	files []string
	err   error
}

// countdownTickMsg is sent on each countdown tick
type countdownTickMsg struct{}

// countdownCompleteMsg is sent when countdown finishes
type countdownCompleteMsg struct{}

// scanLogsDirCmd rescans the logs directory for changes
func scanLogsDirCmd() tea.Cmd {
	return func() tea.Msg {
		files, err := fileops.ScanLogFiles()
		return scanLogsMsg{files: files, err: err}
	}
}

// periodicScanCmd sends a tick every 5 seconds to rescan the logs directory
func periodicScanCmd() tea.Cmd {
	return tea.Every(5*time.Second, func(t time.Time) tea.Msg {
		return scanLogsDirCmd()()
	})
}

// Update handles all the state updates based on incoming messages
func Update(msg tea.Msg, m models.Model) (models.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TermWidth = msg.Width
		m.TermHeight = msg.Height

		// Dynamic left pane width
		targetWidth := m.TermWidth / 3
		minWidth := 25 // Absolute minimum for left pane
		maxWidth := 70 // Sensible maximum for left pane

		if targetWidth < minWidth {
			targetWidth = minWidth
		}
		if targetWidth > maxWidth {
			targetWidth = maxWidth
		}

		// Ensure right pane has at least a minimum width (e.g., 20 chars) if possible
		minRightPaneWidth := 20
		if m.TermWidth-targetWidth < minRightPaneWidth && targetWidth > minWidth {
			targetWidth = m.TermWidth - minRightPaneWidth
			if targetWidth < minWidth { // If terminal is too small for both, prioritize left pane's min
				targetWidth = minWidth
			}
		}
		// Final check: left pane cannot be wider than the terminal itself minus a bit for border/right pane
		if targetWidth >= m.TermWidth-m.LeftPaneStyle.GetHorizontalBorderSize()-5 && m.TermWidth > minWidth+5 {
			targetWidth = m.TermWidth - m.LeftPaneStyle.GetHorizontalBorderSize() - 5
			if targetWidth < minWidth {
				targetWidth = minWidth
			}
		}

		m.LeftPaneWidth = targetWidth
		return m, periodicScanCmd() // Start periodic scanning when window is ready

	case scanLogsMsg:
		if msg.err != nil {
			m.Err = msg.err
			return m, nil
		}

		// Create new menu choices with updated log files
		newChoices := make([]string, 0, len(msg.files)+3) // +3 for Macros, toggle, exit
		newChoices = append(newChoices, msg.files...)
		newChoices = append(newChoices,
			MacrosText, // Add Macros option
			fmt.Sprintf("%s (%s)", CoreProtectToggleBaseText, map[bool]string{true: "ON", false: "OFF"}[m.CoreProtectMode]),
			ExitText,
		)

		// Adjust cursor if needed
		if m.MenuCursor >= len(newChoices) {
			m.MenuCursor = len(newChoices) - 1
		}

		m.MenuChoices = newChoices
		return m, periodicScanCmd() // Continue periodic scanning

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.State {
		case models.MenuView:
			return handleMenuViewInput(msg, m)
		case models.LogView:
			return handleLogViewInput(msg, m)
		case models.SaveInputView:
			return handleSaveInputViewInput(msg, m)
		case models.MacroListView:
			return handleMacroListViewInput(msg, m)
		case models.MacroParameterInputView:
			return handleMacroParameterInputViewInput(msg, m)
		case models.ParameterSelectionView:
			return handleParameterSelectionViewInput(msg, m)
		case models.ParameterValueInputView:
			return handleParameterValueInputViewInput(msg, m)
		case models.CountdownInputView:
			return handleCountdownInputViewInput(msg, m)
		case models.CountdownDisplayView:
			return handleCountdownDisplayViewInput(msg, m)
		}

	case []logparser.LogEntry:
		m.LogEntries = msg
		m.LogCursor = 0
		m.Err = nil
		return m, periodicScanCmd()

	case []coreprotectparser.CoreProtectLogEntry:
		m.CoreProtectLogEntries = msg
		m.LogCursor = 0
		m.Err = nil
		return m, periodicScanCmd()

	case models.SaveSuccessMsg:
		m.SaveMessage = fmt.Sprintf("Logs saved to output/%s", msg.Filename)
		m.State = m.PreviousState
		m.FocusedPane = models.LogFilePane
		m.InputActive = false
		m.SaveFilenameInput = ""
		return m, periodicScanCmd()

	case models.SaveErrorMsg:
		m.SaveMessage = fmt.Sprintf("Error saving: %v", msg.Err)
		m.State = m.PreviousState
		return m, periodicScanCmd()

	case error:
		m.Err = msg
		return m, periodicScanCmd()

	case countdownTickMsg:
		if m.CountdownValue > 0 {
			m.CountdownValue--
			return m, tea.Tick(time.Second, func(time.Time) tea.Msg { return countdownTickMsg{} })
		}
		// Countdown complete - execute macro and return to menu
		if m.SelectedMacroName != "" {
			// For CoreProtect Pager, use completely separate parameter maps
			if m.SelectedMacroName == "CoreProtect Pager" {
				// Create completely new maps for different parameter types
				pagingParams := make(map[string]string)
				lookupParams := make(map[string]string)

				// Get paging parameters directly from input fields, not from MacroParameters
				// This ensures they are clean and not corrupted
				startPage, exists := m.ParameterInputs["startPage"]
				if exists && startPage != "" {
					pagingParams["startPage"] = startPage
				} else {
					pagingParams["startPage"] = "1" // Default
				}

				endPage, exists := m.ParameterInputs["endPage"]
				if exists && endPage != "" {
					pagingParams["endPage"] = endPage
				} else {
					pagingParams["endPage"] = "5" // Default
				}

				delayMs, exists := m.ParameterInputs["delayMs"]
				if exists && delayMs != "" {
					pagingParams["delayMs"] = delayMs
				} else {
					pagingParams["delayMs"] = "500" // Default
				}

				// Get lookup parameters from ParameterValues, not from MacroParameters
				if users, ok := m.ParameterValues["users"]; ok && len(users) > 0 {
					lookupParams["users"] = strings.Join(users, ",")
				}

				if actions, ok := m.ParameterValues["actions"]; ok && len(actions) > 0 {
					lookupParams["actions"] = strings.Join(actions, ",")
				}

				// Combine parameters only at the last moment
				finalParams := make(map[string]string)

				// Add paging parameters
				finalParams["startPage"] = pagingParams["startPage"]
				finalParams["endPage"] = pagingParams["endPage"]
				finalParams["delayMs"] = pagingParams["delayMs"]

				// Add lookup parameters
				if users, ok := lookupParams["users"]; ok {
					finalParams["users"] = users
				}

				if actions, ok := lookupParams["actions"]; ok {
					finalParams["actions"] = actions
				}

				// Get radius from ParameterValues if available
				if radiusValues, ok := m.ParameterValues["radius"]; ok && len(radiusValues) > 0 {
					finalParams["radius"] = radiusValues[0]
				}

				// Get time from ParameterValues if available
				if timeValues, ok := m.ParameterValues["time"]; ok && len(timeValues) > 0 {
					finalParams["time"] = timeValues[0]
				}

				// Execute with final parameters
				err := macros.ExecuteMacro(m.SelectedMacroName, finalParams)
				if err != nil {
					m.Err = err
				}
			} else {
				// For other macros, execute normally
				err := macros.ExecuteMacro(m.SelectedMacroName, m.MacroParameters)
				if err != nil {
					m.Err = err
				}
			}

			m.SelectedMacroName = ""
			m.MacroParameters = make(map[string]string)   // Clear parameters
			m.ParameterInputs = make(map[string]string)   // Clear parameter inputs
			m.ParameterValues = make(map[string][]string) // Clear parameter values
		}
		m.State = models.MenuView
		m.FocusedPane = models.LogFilePane
		m.InputActive = false
		m.CountdownActive = false
		return m, nil

	case countdownCompleteMsg:
		m.State = models.MenuView
		m.FocusedPane = models.LogFilePane
		m.InputActive = false
		m.CountdownActive = false
		m.CountdownValue = 0
		// Execute the selected macro
		if m.SelectedMacroName != "" {
			go func() {
				// In a real implementation, this would execute the macro
				// macros.ExecuteMacro(m.SelectedMacroName)
			}()
			m.SelectedMacroName = ""
		}
		return m, nil
	}

	return m, cmd
}

// handleMenuViewInput handles input when in menu view
func handleMenuViewInput(msg tea.KeyMsg, m models.Model) (models.Model, tea.Cmd) {
	if m.FocusedPane == models.LogFilePane {
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.MenuCursor > 0 {
				m.MenuCursor--
			}
		case "down", "j":
			if m.MenuCursor < len(m.MenuChoices)-1 {
				m.MenuCursor++
			}
		case "tab":
			if !m.CoreProtectMode {
				m.FocusedPane = models.FilterPane
				m.InputActive = true
			}
		case "enter":
			selectedChoice := m.MenuChoices[m.MenuCursor]
			if selectedChoice == ExitText {
				return m, tea.Quit
			} else if selectedChoice == MacrosText {
				m.State = models.MacroListView
				m.FocusedPane = models.MacroListPane // Focus the macro list in the right pane
				m.MacroCursor = 0                    // Reset macro cursor
				m.SaveMessage = ""                   // Clear any previous save message
				m.Err = nil                          // Clear any previous error
				return m, nil
			} else if strings.HasPrefix(selectedChoice, CoreProtectToggleBaseText) {
				m.CoreProtectMode = !m.CoreProtectMode
				// Update the menu choice text directly
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
				return m, nil
			} else { // It's a log file
				m.State = models.LogView
				m.LogEntries = []logparser.LogEntry{}
				m.CoreProtectLogEntries = []coreprotectparser.CoreProtectLogEntry{}
				m.LogCursor = 0
				m.Err = nil
				m.SaveMessage = ""
				return m, loadLogFileCmd(selectedChoice, m.Filters, m.CoreProtectMode)
			}
		}
	} else if m.FocusedPane == models.FilterPane { // Input for filter
		switch msg.String() {
		case "q":
			if m.FilterInput == "" {
				return m, tea.Quit
			}
			m.FilterInput += "q"
		case "tab":
			m.FocusedPane = models.LogFilePane
			m.InputActive = false
		case "esc":
			m.FocusedPane = models.LogFilePane
			m.InputActive = false
		case "enter":
			if m.FilterInput != "" {
				m.Filters = append(m.Filters, m.FilterInput)
				m.FilterInput = ""
				if m.State == models.LogView && !m.CoreProtectMode {
					currentLogFile := m.MenuChoices[m.MenuCursor]
					if !strings.HasPrefix(currentLogFile, CoreProtectToggleBaseText) && currentLogFile != ExitText {
						return m, loadLogFileCmd(currentLogFile, m.Filters, m.CoreProtectMode)
					}
				}
			}
		case "backspace":
			if len(m.FilterInput) > 0 {
				m.FilterInput = m.FilterInput[:len(m.FilterInput)-1]
			}
		default:
			if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
				m.FilterInput += string(msg.Runes)
			}
		}
	}
	return m, nil
}

// handleLogViewInput handles input when in log view
func handleLogViewInput(msg tea.KeyMsg, m models.Model) (models.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.FocusedPane != models.FilterPane && m.LogCursor > 0 {
			m.LogCursor--
		}
	case "down", "j":
		currentLogListSize := 0
		if m.CoreProtectMode {
			currentLogListSize = len(m.CoreProtectLogEntries)
		} else {
			currentLogListSize = len(m.LogEntries)
		}

		if m.FocusedPane != models.FilterPane && m.LogCursor < currentLogListSize-1 {
			m.LogCursor++
		}
	case "e":
		if (m.CoreProtectMode && len(m.CoreProtectLogEntries) > 0) || (!m.CoreProtectMode && len(m.LogEntries) > 0) {
			m.PreviousState = m.State
			m.State = models.SaveInputView
			m.SaveFilenameInput = ""
			m.SaveMessage = ""
		}
	case "esc":
		m.State = models.MenuView
		m.FocusedPane = models.LogFilePane
		m.InputActive = false
		m.SaveMessage = ""
	case "tab":
		if !m.CoreProtectMode {
			if m.FocusedPane == models.LogFilePane {
				m.FocusedPane = models.FilterPane
				m.InputActive = true
			} else {
				m.FocusedPane = models.LogFilePane
				m.InputActive = false
			}
		} else {
			m.FocusedPane = models.LogFilePane
			m.InputActive = false
		}
	}
	// If in LogView and focused on LogFilePane (left pane, but not filter input),
	// TAB should still switch to filter input if available.
	if m.State == models.LogView && m.FocusedPane == models.LogFilePane && !m.CoreProtectMode && msg.String() == "tab" {
		m.FocusedPane = models.FilterPane
		m.InputActive = true
	}

	return m, nil
}

// handleMacroListViewInput handles input when in macro list view
func handleMacroListViewInput(msg tea.KeyMsg, m models.Model) (models.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		// Currently, 'q' does nothing here. 'esc' is used to go back.
	case "up", "k":
		if m.FocusedPane == models.MacroListPane && m.MacroCursor > 0 {
			m.MacroCursor--
		}
	case "down", "j":
		if m.FocusedPane == models.MacroListPane && m.MacroCursor < len(m.MacroChoices)-1 {
			m.MacroCursor++
		}
	case "enter":
		if m.FocusedPane == models.MacroListPane && len(m.MacroChoices) > 0 && m.MacroCursor < len(m.MacroChoices) {
			m.SelectedMacroName = m.MacroChoices[m.MacroCursor]

			// Get the selected macro to check if it has parameters
			var selectedMacro *macros.Macro
			for _, macro := range macros.GetAvailableMacros() {
				if macro.Name == m.SelectedMacroName {
					tempMacro := macro
					selectedMacro = &tempMacro
					break
				}
			}

			// Clear any existing parameters
			m.MacroParameters = make(map[string]string)
			m.ParameterInputs = make(map[string]string)

			// If the macro has parameters, show parameter input view
			if selectedMacro != nil && len(selectedMacro.Parameters) > 0 {
				m.PreviousState = m.State
				m.State = models.MacroParameterInputView
				m.ParameterCursor = 0
				m.InputActive = true

				// Initialize parameter inputs with default values
				m.ParameterInputs = make(map[string]string)
				for _, param := range selectedMacro.Parameters {
					m.ParameterInputs[param.Name] = param.DefaultValue
				}

				// For CoreProtect Pager, initialize additional fields
				if selectedMacro.Name == "CoreProtect Pager" {
					// Initialize available parameters for selection with all possible action parameters
					m.AvailableParameters = []string{
						"users",
						"action:+block",
						"action:+container",
						"action:+inventory",
						"action:+session",
						"action:-block",
						"action:-container",
						"action:-inventory",
						"action:-item",
						"action:-session",
						"action:block",
						"action:chat",
						"action:click",
						"action:command",
						"action:container",
						"action:inventory",
						"action:item",
						"action:kill",
						"action:session",
						"action:sign",
						"action:username",
						"radius",
						"time",
					}

					// Initialize parameter values map for multi-value parameters
					m.ParameterValues = make(map[string][]string)

					// Initialize MacroParameters with default values
					m.MacroParameters = make(map[string]string)
					for _, param := range selectedMacro.Parameters {
						m.MacroParameters[param.Name] = param.DefaultValue
					}
				}
			} else {
				// No parameters, go straight to countdown
				m.PreviousState = m.State
				m.State = models.CountdownInputView
				m.CountdownInput = ""
				m.CountdownMessage = ""
				m.InputActive = true
			}
		}
	case "esc":
		m.State = models.MenuView
		m.FocusedPane = models.LogFilePane
		m.InputActive = false
		m.SaveMessage = ""
	case "tab":
		m.FocusedPane = models.LogFilePane
		m.InputActive = false
	}
	return m, nil
}

// handleSaveInputViewInput handles input when in save input view
func handleSaveInputViewInput(msg tea.KeyMsg, m models.Model) (models.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.State = m.PreviousState
		m.SaveFilenameInput = ""
		m.SaveMessage = ""
		m.FocusedPane = models.LogFilePane
		m.InputActive = false
	case "enter":
		if m.SaveFilenameInput != "" {
			return m, saveFileCmd(m)
		} else {
			m.SaveMessage = "Filename cannot be empty."
		}
	case "backspace":
		if len(m.SaveFilenameInput) > 0 {
			m.SaveFilenameInput = m.SaveFilenameInput[:len(m.SaveFilenameInput)-1]
		}
	default:
		if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
			m.SaveFilenameInput += string(msg.Runes)
		}
	}
	return m, nil
}

// saveFileCmd creates a command to save the log entries
func saveFileCmd(m models.Model) tea.Cmd {
	return func() tea.Msg {
		if m.SaveFilenameInput == "" {
			return models.SaveErrorMsg{Err: fmt.Errorf("filename cannot be empty")}
		}
		var err error
		if m.CoreProtectMode {
			err = fileops.SaveCoreProtectLogsToFile(m.CoreProtectLogEntries, m.SaveFilenameInput)
		} else {
			err = fileops.SaveStandardLogsToFile(m.LogEntries, m.SaveFilenameInput)
		}
		if err != nil {
			return models.SaveErrorMsg{Err: err}
		}
		return models.SaveSuccessMsg{Filename: m.SaveFilenameInput}
	}
}

// handleMacroParameterInputViewInput handles input when in parameter input view
func handleMacroParameterInputViewInput(msg tea.KeyMsg, m models.Model) (models.Model, tea.Cmd) {
	// Get the selected macro to access its parameters
	var selectedMacro *macros.Macro
	for _, macro := range macros.GetAvailableMacros() {
		if macro.Name == m.SelectedMacroName {
			tempMacro := macro
			selectedMacro = &tempMacro
			break
		}
	}

	if selectedMacro == nil {
		// Something went wrong, go back to macro list
		m.State = models.MacroListView
		return m, nil
	}

	// Get the current parameter
	if len(selectedMacro.Parameters) == 0 {
		// No parameters, go to countdown
		m.State = models.CountdownInputView
		m.CountdownInput = ""
		m.CountdownMessage = ""
		return m, nil
	}

	// Special handling for CoreProtect Pager
	if selectedMacro.Name == "CoreProtect Pager" {
		// Check if cursor is on the "Add Parameter" button
		if m.ParameterCursor == len(selectedMacro.Parameters) {
			switch msg.String() {
			case "enter":
				// Go to parameter selection view
				m.State = models.ParameterSelectionView
				m.ParameterCursor = 0
				return m, nil
			case "tab":
				// Wrap around to first parameter
				m.ParameterCursor = 0
				return m, nil
			case "shift+tab":
				// Move to last parameter
				m.ParameterCursor = len(selectedMacro.Parameters) - 1
				return m, nil
			}
		}
	}

	// Make sure cursor is in valid range
	if m.ParameterCursor >= len(selectedMacro.Parameters) {
		m.ParameterCursor = len(selectedMacro.Parameters) - 1
	}

	currentParam := selectedMacro.Parameters[m.ParameterCursor]

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "ctrl+e":
		// If it's CoreProtect Pager, go to parameter selection view
		if selectedMacro.Name == "CoreProtect Pager" {
			m.State = models.ParameterSelectionView
			m.ParameterCursor = 0
			return m, nil
		}
	case "esc":
		// Cancel parameter input and go back to macro list
		m.State = models.MacroListView
		m.ParameterInputs = make(map[string]string)
		m.MacroParameters = make(map[string]string)
		m.ParameterMessage = ""
		m.ParameterCursor = 0
	case "tab":
		// Move to next parameter field
		if m.ParameterCursor < len(selectedMacro.Parameters)-1 {
			m.ParameterCursor++
		} else if selectedMacro.Name == "CoreProtect Pager" {
			// For CoreProtect Pager, move to the "Add Parameter" button
			m.ParameterCursor = len(selectedMacro.Parameters)
		} else {
			m.ParameterCursor = 0 // Wrap around to first parameter
		}
	case "shift+tab":
		// Move to previous parameter field
		if m.ParameterCursor > 0 {
			m.ParameterCursor--
		} else if selectedMacro.Name == "CoreProtect Pager" {
			// For CoreProtect Pager, move to the "Add Parameter" button
			m.ParameterCursor = len(selectedMacro.Parameters)
		} else {
			m.ParameterCursor = len(selectedMacro.Parameters) - 1 // Wrap around to last parameter
		}
	case "enter":
		// Save parameters and move to countdown
		m.MacroParameters = make(map[string]string)

		// For CoreProtect Pager, handle parameters with complete separation
		if selectedMacro.Name == "CoreProtect Pager" {
			// ONLY save the basic paging parameters from the input fields
			// These are used for the "/co page" commands
			startPage, exists := m.ParameterInputs["startPage"]
			if exists && startPage != "" {
				m.MacroParameters["startPage"] = startPage
			} else {
				m.MacroParameters["startPage"] = "1" // Default
			}

			endPage, exists := m.ParameterInputs["endPage"]
			if exists && endPage != "" {
				m.MacroParameters["endPage"] = endPage
			} else {
				m.MacroParameters["endPage"] = "5" // Default
			}

			delayMs, exists := m.ParameterInputs["delayMs"]
			if exists && delayMs != "" {
				m.MacroParameters["delayMs"] = delayMs
			} else {
				m.MacroParameters["delayMs"] = "500" // Default
			}

			// SEPARATELY handle the lookup parameters
			// These are used for the "/co lookup" command

			// Users
			if users, ok := m.ParameterValues["users"]; ok && len(users) > 0 {
				m.MacroParameters["users"] = strings.Join(users, ",")
			}

			// Actions
			if actions, ok := m.ParameterValues["actions"]; ok && len(actions) > 0 {
				m.MacroParameters["actions"] = strings.Join(actions, ",")
			}

			// Radius (single value parameter)
			if m.MacroParameters["radius"] == "" { // Only set if not already set
				if radius, ok := m.ParameterInputs["radius"]; ok && radius != "" {
					m.MacroParameters["radius"] = radius
				}
			}

			// Time (single value parameter)
			if m.MacroParameters["time"] == "" { // Only set if not already set
				if timeParam, ok := m.ParameterInputs["time"]; ok && timeParam != "" {
					m.MacroParameters["time"] = timeParam
				}
			}
		} else {
			// For other macros, handle parameters normally
			for _, param := range selectedMacro.Parameters {
				value, exists := m.ParameterInputs[param.Name]
				if exists {
					m.MacroParameters[param.Name] = value
				} else {
					m.MacroParameters[param.Name] = param.DefaultValue
				}
			}
		}

		// Move to countdown input
		m.State = models.CountdownInputView
		m.CountdownInput = ""
		m.CountdownMessage = ""
		m.InputActive = true
	case "backspace":
		// Delete character from current parameter input
		currentValue := m.ParameterInputs[currentParam.Name]
		if len(currentValue) > 0 {
			m.ParameterInputs[currentParam.Name] = currentValue[:len(currentValue)-1]
		}
	default:
		// Add character to current parameter input
		if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
			if m.ParameterInputs == nil {
				m.ParameterInputs = make(map[string]string)
			}
			currentValue, exists := m.ParameterInputs[currentParam.Name]
			if !exists {
				currentValue = ""
			}
			m.ParameterInputs[currentParam.Name] = currentValue + string(msg.Runes)
		}
	}
	return m, nil
}

// handleParameterSelectionViewInput handles input when in parameter selection view
func handleParameterSelectionViewInput(msg tea.KeyMsg, m models.Model) (models.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Go back to parameter input view
		m.State = models.MacroParameterInputView
		return m, nil
	case "ctrl+d":
		// Done with parameters, go to countdown
		m.State = models.CountdownInputView
		m.CountdownInput = ""
		m.CountdownMessage = ""
		m.InputActive = true
		return m, nil
	case "up", "k":
		// Move cursor up
		if m.ParameterCursor > 0 {
			m.ParameterCursor--
		}
	case "down", "j":
		// Move cursor down
		if m.ParameterCursor < len(m.AvailableParameters)-1 {
			m.ParameterCursor++
		}
	case "enter":
		// Select parameter
		if m.ParameterCursor < len(m.AvailableParameters) {
			selectedParam := m.AvailableParameters[m.ParameterCursor]

			// Check if it's an action parameter
			if strings.HasPrefix(selectedParam, "action:") {
				// For action parameters, add them directly to the active parameters
				actionType := strings.TrimPrefix(selectedParam, "action:")

				// Initialize the actions slice if needed
				if m.ParameterValues["actions"] == nil {
					m.ParameterValues["actions"] = []string{}
				}

				// Add the action to the list
				m.ParameterValues["actions"] = append(m.ParameterValues["actions"], actionType)

				// DO NOT update MacroParameters at all
				// We'll only combine parameters at the very end when executing the macro
				// This ensures complete separation and prevents any corruption

				// Set a success message to indicate the parameter was added (will be displayed in green)
				m.ParameterMessage = fmt.Sprintf("SUCCESS: Action parameter '%s' added", actionType)

				// Return to the parameter input view
				m.State = models.MacroParameterInputView
				return m, nil
			} else {
				// For other parameters (users, radius, time)
				m.SelectedParameter = selectedParam

				// Initialize parameter value input
				m.ParameterValueInput = ""
				m.ParameterValueMessage = ""

				// Initialize parameter values map if needed
				if m.ParameterValues == nil {
					m.ParameterValues = make(map[string][]string)
				}

				// Go to parameter value input view
				m.State = models.ParameterValueInputView
			}
		}
	}
	return m, nil
}

// handleParameterValueInputViewInput handles input when in parameter value input view
func handleParameterValueInputViewInput(msg tea.KeyMsg, m models.Model) (models.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Go back to parameter selection view
		m.State = models.ParameterSelectionView
		return m, nil
	case "enter":
		// Add parameter value
		if m.ParameterValueInput != "" {
			// For user parameter
			if m.SelectedParameter == "users" {
				// Initialize the slice if needed
				if m.ParameterValues["users"] == nil {
					m.ParameterValues["users"] = []string{}
				}

				// Add the value to the slice
				m.ParameterValues["users"] = append(
					m.ParameterValues["users"],
					m.ParameterValueInput,
				)

				// DO NOT update MacroParameters at all
				// We'll only combine parameters at the very end when executing the macro
				// This ensures complete separation and prevents any corruption

				// Clear the input for another value
				m.ParameterValueInput = ""
				m.ParameterValueMessage = "User added. Enter another or press ESC to go back."
			} else {
				// For single-value parameters (radius, time)
				if m.SelectedParameter == "radius" || m.SelectedParameter == "time" {
					// Store in a separate field to avoid mixing with paging parameters
					if m.ParameterValues[m.SelectedParameter] == nil {
						m.ParameterValues[m.SelectedParameter] = []string{}
					}
					m.ParameterValues[m.SelectedParameter] = []string{m.ParameterValueInput}

					// DO NOT update MacroParameters at all
					// We'll only combine parameters at the very end when executing the macro
					// This ensures complete separation and prevents any corruption

					// Set a success message to indicate the parameter was added
					m.ParameterMessage = fmt.Sprintf("SUCCESS: %s parameter set to '%s'",
						m.SelectedParameter, m.ParameterValueInput)
				} else {
					// For other parameters (not lookup related)
					// Store directly in MacroParameters
					if m.MacroParameters == nil {
						m.MacroParameters = make(map[string]string)
					}
					m.MacroParameters[m.SelectedParameter] = m.ParameterValueInput
				}

				// Go back to parameter selection
				m.State = models.ParameterSelectionView
			}
		} else {
			m.ParameterValueMessage = "Value cannot be empty"
		}
	case "backspace":
		// Delete character from input
		if len(m.ParameterValueInput) > 0 {
			m.ParameterValueInput = m.ParameterValueInput[:len(m.ParameterValueInput)-1]
		}
	default:
		// Add character to input
		if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
			m.ParameterValueInput += string(msg.Runes)
		}
	}
	return m, nil
}

// handleCountdownInputViewInput handles input when in countdown input view
func handleCountdownInputViewInput(msg tea.KeyMsg, m models.Model) (models.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.State = m.PreviousState
		m.CountdownInput = ""
		m.CountdownMessage = ""
		m.FocusedPane = models.LogFilePane
		m.InputActive = false
	case "enter":
		if m.CountdownInput != "" {
			// Validate countdown input is a positive integer
			var seconds int
			_, err := fmt.Sscanf(m.CountdownInput, "%d", &seconds)
			if err != nil || seconds <= 0 {
				m.CountdownMessage = "Invalid time - must be positive number"
				return m, nil
			}

			m.CountdownValue = seconds
			m.CountdownActive = true
			m.State = models.CountdownDisplayView
			m.CountdownInput = ""
			m.CountdownMessage = ""
			return m, tea.Tick(time.Second, func(time.Time) tea.Msg { return countdownTickMsg{} })
		} else {
			m.CountdownMessage = "Time cannot be empty"
		}
	case "backspace":
		if len(m.CountdownInput) > 0 {
			m.CountdownInput = m.CountdownInput[:len(m.CountdownInput)-1]
		}
	default:
		if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
			m.CountdownInput += string(msg.Runes)
		}
	}
	return m, nil
}

// handleCountdownDisplayViewInput handles input when in countdown display view
func handleCountdownDisplayViewInput(msg tea.KeyMsg, m models.Model) (models.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.State = models.MenuView
		m.FocusedPane = models.LogFilePane
		m.InputActive = false
		m.CountdownActive = false
		m.CountdownValue = 0
	}
	return m, nil
}

// loadLogFileCmd is a command that sends the loaded entries back as a message
func loadLogFileCmd(filePath string, filters []string, coreProtectMode bool) tea.Cmd {
	return func() tea.Msg {
		if coreProtectMode {
			// For CoreProtect, we read the whole file content then parse
			content, err := fileops.ReadFileContent(filePath)
			if err != nil {
				return fmt.Errorf("failed to read CoreProtect log file %s: %w", filePath, err)
			}
			cpLog, err := coreprotectparser.ParseLogContent(content)
			if err != nil {
				return fmt.Errorf("failed to parse CoreProtect log file %s: %w", filePath, err)
			}
			return cpLog.Entries
		} else {
			// Standard log parsing
			parser, err := logparser.NewParser()
			if err != nil {
				return err
			}

			// Use our gzip-aware reader for all files
			content, err := fileops.ReadFileContent(filePath)
			if err != nil {
				return fmt.Errorf("failed to read log file %s: %w", filePath, err)
			}

			entries, err := parser.ParseContent(content, filters)
			if err != nil {
				return fmt.Errorf("failed to parse log content from %s: %w", filePath, err)
			}
			return entries
		}
	}
}
