package ui

import (
	"fmt"
	"strings"
	"time"

	"goparselogs/internal/fileops"
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
		newChoices := make([]string, 0, len(msg.files)+2)
		newChoices = append(newChoices, msg.files...)
		newChoices = append(newChoices,
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
			} else if strings.HasPrefix(selectedChoice, CoreProtectToggleBaseText) {
				m.CoreProtectMode = !m.CoreProtectMode
				return m, nil
			} else {
				m.State = models.LogView
				m.LogEntries = []logparser.LogEntry{}
				m.CoreProtectLogEntries = []coreprotectparser.CoreProtectLogEntry{}
				m.LogCursor = 0
				m.Err = nil
				return m, loadLogFileCmd(selectedChoice, m.Filters, m.CoreProtectMode)
			}
		}
	} else if m.FocusedPane == models.FilterPane {
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
