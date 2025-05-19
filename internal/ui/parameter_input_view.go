package ui

import (
	"fmt"
	"strings"

	"goparselogs/internal/macros"
	"goparselogs/internal/models"

	"github.com/charmbracelet/lipgloss"
)

// renderMacroParameterInputView renders the parameter input view as a full-screen modal
// showing all parameters at once with the ability to tab through them
func renderMacroParameterInputView(m models.Model, width int) string {
	var view strings.Builder

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
		view.WriteString("Error: Selected macro not found.")
		return view.String()
	}

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

	// Description
	modalContent.WriteString(fmt.Sprintf("Description: %s\n\n", selectedMacro.Description))

	// No parameters case
	if len(selectedMacro.Parameters) == 0 {
		modalContent.WriteString("This macro has no parameters.\n\n")
	} else {
		// Show all parameters as a form
		modalContent.WriteString("Parameters:\n\n")

		// Calculate max parameter name length for alignment
		maxNameLength := 0
		for _, param := range selectedMacro.Parameters {
			if len(param.Name) > maxNameLength {
				maxNameLength = len(param.Name)
			}
		}

		// Render each parameter field
		for i, param := range selectedMacro.Parameters {
			// Get the current input value for this parameter
			inputValue, exists := m.ParameterInputs[param.Name]
			if !exists {
				// If not set yet, use default value
				inputValue = param.DefaultValue
				// Store it in the inputs map
				if m.ParameterInputs == nil {
					m.ParameterInputs = make(map[string]string)
				}
				m.ParameterInputs[param.Name] = inputValue
			}

			// Parameter name with padding for alignment
			paramNameStyle := lipgloss.NewStyle().Bold(true)
			if i == m.ParameterCursor {
				paramNameStyle = paramNameStyle.Foreground(lipgloss.Color("63"))
			}

			paddedName := param.Name
			if len(paddedName) < maxNameLength {
				paddedName = paddedName + strings.Repeat(" ", maxNameLength-len(paddedName))
			}

			modalContent.WriteString(paramNameStyle.Render(paddedName) + ": ")

			// Parameter description
			modalContent.WriteString(param.Description + "\n")

			// Input field
			inputStyle := m.InputStyle
			if i == m.ParameterCursor && m.InputActive {
				inputStyle = m.FocusedInputStyle
			}

			// Limit input width to modal width
			maxInputWidth := modalWidth - 10
			if maxInputWidth < 10 {
				maxInputWidth = 10
			}
			inputStyle = inputStyle.Copy().Width(maxInputWidth)

			// Ensure proper alignment of the input field box
			modalContent.WriteString(inputStyle.Render(inputValue) + "\n\n")
		}
	}

	// Show any validation message
	if m.ParameterMessage != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
		modalContent.WriteString(errorStyle.Render(m.ParameterMessage) + "\n\n")
	}

	// Show help text
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	modalContent.WriteString("\n" + helpStyle.Render("TAB: Next field • SHIFT+TAB: Previous field • ENTER: Confirm • ESC: Cancel"))

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
