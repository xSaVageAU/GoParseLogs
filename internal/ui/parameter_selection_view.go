package ui

import (
	"fmt"
	"strings"

	"goparselogs/internal/models"

	"github.com/charmbracelet/lipgloss"
)

// renderParameterSelectionView renders the view for selecting which parameter to add
func renderParameterSelectionView(m models.Model, width int) string {
	var view strings.Builder

	// Create a centered modal style
	modalWidth := Min(width-10, 80) // Max width of 80 or screen width minus padding
	if modalWidth < 40 {
		modalWidth = Min(width, 40) // Minimum width of 40
	}

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Width(modalWidth).
		Align(lipgloss.Center)

	var modalContent strings.Builder

	// Title with highlight
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	modalContent.WriteString(titleStyle.Render(fmt.Sprintf("Configure Macro: %s\n\n", m.SelectedMacroName)))

	// Show available parameters to add
	modalContent.WriteString("Select Parameter to Add:\n")
	for i, param := range m.AvailableParameters {
		cursor := "  "
		if i == m.ParameterCursor {
			cursor = "> "
			param = m.HighlightStyle.Render(param)
		}
		modalContent.WriteString(fmt.Sprintf("%s%s\n", cursor, param))
	}

	// Show help text
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	modalContent.WriteString("\n" + helpStyle.Render("UP/DOWN: Navigate • ENTER: Select • ESC: Back • CTRL+D: Done"))

	// Render the modal and center it on screen
	modal := modalStyle.Render(modalContent.String())

	// Calculate vertical position to center the modal
	lines := strings.Count(modal, "\n") + 1
	verticalPadding := (m.TermHeight - lines) / 2
	if verticalPadding < 0 {
		verticalPadding = 0
	}

	// Add vertical padding
	for i := 0; i < verticalPadding; i++ {
		view.WriteString("\n")
	}

	// Center the modal horizontally
	centeredStyle := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	view.WriteString(centeredStyle.Render(modal))

	return view.String()
}

// renderParameterValueInputView renders the view for entering a value for a selected parameter
func renderParameterValueInputView(m models.Model, width int) string {
	var view strings.Builder

	// Create a centered modal style
	modalWidth := Min(width-10, 80) // Max width of 80 or screen width minus padding
	if modalWidth < 40 {
		modalWidth = Min(width, 40) // Minimum width of 40
	}

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Width(modalWidth).
		Align(lipgloss.Center)

	var modalContent strings.Builder

	// Title with highlight
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))

	// Parameter description
	var description string
	switch m.SelectedParameter {
	case "users":
		description = "Enter a username to lookup (e.g., 'player1')"
		titleStyle = titleStyle.Copy().Foreground(lipgloss.Color("10")) // Green for user
		modalContent.WriteString(titleStyle.Render("Enter Username\n\n"))
	case "actions":
		description = "Enter an action to lookup (e.g., 'block', 'chat', 'command')"
		titleStyle = titleStyle.Copy().Foreground(lipgloss.Color("11")) // Yellow for action
		modalContent.WriteString(titleStyle.Render("Enter Action\n\n"))
	case "radius":
		description = "Enter a radius for the lookup (e.g., '10')"
		modalContent.WriteString(titleStyle.Render("Enter Radius\n\n"))
	case "time":
		description = "Enter a time parameter (e.g., '1d' for 1 day, '12h' for 12 hours)"
		modalContent.WriteString(titleStyle.Render("Enter Time Parameter\n\n"))
	default:
		description = "Enter a value for " + m.SelectedParameter
		modalContent.WriteString(titleStyle.Render(fmt.Sprintf("Enter %s\n\n", m.SelectedParameter)))
	}
	modalContent.WriteString(description + "\n\n")

	// Show current values for this parameter if it's a multi-value parameter
	if m.SelectedParameter == "users" || m.SelectedParameter == "actions" {
		values := m.ParameterValues[m.SelectedParameter]
		if len(values) > 0 {
			modalContent.WriteString("Current values:\n")
			for _, value := range values {
				modalContent.WriteString("  - " + value + "\n")
			}
			modalContent.WriteString("\n")
		}
	}

	// Input field
	inputStyle := m.FocusedInputStyle
	inputStyle = inputStyle.Copy().Width(modalWidth - 10)
	modalContent.WriteString(inputStyle.Render(m.ParameterValueInput) + "\n\n")

	// Show any validation message
	if m.ParameterValueMessage != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
		modalContent.WriteString(errorStyle.Render(m.ParameterValueMessage) + "\n\n")
	}

	// Show help text
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	modalContent.WriteString("\n" + helpStyle.Render("ENTER: Add Value • ESC: Back"))

	// Render the modal and center it on screen
	modal := modalStyle.Render(modalContent.String())

	// Calculate vertical position to center the modal
	lines := strings.Count(modal, "\n") + 1
	verticalPadding := (m.TermHeight - lines) / 2
	if verticalPadding < 0 {
		verticalPadding = 0
	}

	// Add vertical padding
	for i := 0; i < verticalPadding; i++ {
		view.WriteString("\n")
	}

	// Center the modal horizontally
	centeredStyle := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	view.WriteString(centeredStyle.Render(modal))

	return view.String()
}
