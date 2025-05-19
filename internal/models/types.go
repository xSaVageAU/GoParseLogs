package models

import (
	"goparselogs/pkg/coreprotectparser"
	"goparselogs/pkg/logparser"

	"github.com/charmbracelet/lipgloss"
)

type AppState int

const (
	MenuView      AppState = iota // Main menu with log files and filter input
	LogView                       // View for displaying logs
	SaveInputView                 // View for entering filename to save
	MacroListView                 // View for displaying list of macros
)

// Messages for save operation
type SaveSuccessMsg struct{ Filename string }
type SaveErrorMsg struct{ Err error }

type FocusablePane int

const (
	LogFilePane   FocusablePane = iota // Focus on the log file list in the left pane
	FilterPane                         // Focus on the filter input in the left pane
	MacroListPane                      // Focus on the macro list in the right pane (when State is MacroListView)
)

type Model struct {
	State           AppState
	FocusedPane     FocusablePane // To manage focus: LogFilePane/FilterPane (left), MacroListPane (right)
	PreviousState   AppState      // To store the state before entering saveInputView
	CoreProtectMode bool          // True if CoreProtect parsing is enabled

	// Window / Layout
	TermWidth     int
	TermHeight    int
	LeftPaneWidth int // Desired width for the left (menu) pane

	// Menu View / Shared
	MenuChoices []string // Log files, "Macros", "Exit"
	MenuCursor  int      // For LogFilePane (selection in MenuChoices)
	FilterInput string   // Current text in filter input field
	Filters     []string // List of active filters
	InputActive bool     // True when filterInput has focus (i.e., focusedPane == filterPane)

	// Log View
	LogEntries            []logparser.LogEntry
	CoreProtectLogEntries []coreprotectparser.CoreProtectLogEntry
	LogCursor             int   // cursor for log view (applies to either type of log)
	Err                   error // General errors

	// Save Input View
	SaveFilenameInput string
	SaveMessage       string // To display "Saved!" or "Error saving."

	// Macro List View (Right Pane when "Macros" is selected)
	MacroChoices []string // List of available macros
	MacroCursor  int      // Cursor for MacroChoices

	// Styles
	HighlightStyle    lipgloss.Style
	SubtleStyle       lipgloss.Style
	InputStyle        lipgloss.Style
	FocusedInputStyle lipgloss.Style
	LeftPaneStyle     lipgloss.Style
	RightPaneStyle    lipgloss.Style
	ErrorStyle        lipgloss.Style
	SuccessStyle      lipgloss.Style
}
