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
// For CoreProtect Pager, it also shows active parameters and a button to add more parameters
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

	// Description - more compact
	modalContent.WriteString(fmt.Sprintf("%s\n", selectedMacro.Description))

	// Special handling for CoreProtect Pager
	if selectedMacro.Name == "CoreProtect Pager" {
		// Show paging parameters
		modalContent.WriteString("Paging Parameters:\n")
	}

	// No parameters case
	if len(selectedMacro.Parameters) == 0 {
		modalContent.WriteString("This macro has no parameters.\n\n")
	} else {
		// Show all parameters as a form
		if selectedMacro.Name != "CoreProtect Pager" {
			modalContent.WriteString("Parameters:\n\n")
		}

		// Calculate max parameter name length for alignment
		maxNameLength := 0
		var coreProtectParams []macros.MacroParameter

		// For CoreProtect Pager, only show the 3 required fields
		if selectedMacro.Name == "CoreProtect Pager" {
			// Only include startPage, endPage, and delayMs
			for _, param := range selectedMacro.Parameters {
				if param.Name == "startPage" || param.Name == "endPage" || param.Name == "delayMs" {
					coreProtectParams = append(coreProtectParams, param)
					if len(param.Name) > maxNameLength {
						maxNameLength = len(param.Name)
					}
				}
			}
		} else {
			// For other macros, show all parameters
			coreProtectParams = selectedMacro.Parameters
			for _, param := range selectedMacro.Parameters {
				if len(param.Name) > maxNameLength {
					maxNameLength = len(param.Name)
				}
			}
		}

		// Render each parameter field
		for i, param := range coreProtectParams {
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

			// Parameter name on its own line
			modalContent.WriteString(paramNameStyle.Render(paddedName) + ":\n")

			// Input field on next line
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

	// Show any validation or success message
	if m.ParameterMessage != "" {
		if strings.HasPrefix(m.ParameterMessage, "SUCCESS:") {
			// Success message in green
			successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
			modalContent.WriteString(successStyle.Render(m.ParameterMessage) + "\n\n")
		} else {
			// Error message in red
			errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
			modalContent.WriteString(errorStyle.Render(m.ParameterMessage) + "\n\n")
		}
	}

	// Show active parameters list if it's CoreProtect Pager
	if selectedMacro.Name == "CoreProtect Pager" {
		// Create a more compact box for active parameters
		activeParamsStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1).
			MarginTop(1)

		var activeParamsContent strings.Builder
		activeParamsContent.WriteString("Active Parameters:\n")

		// Check if we have any lookup parameters
		hasLookupParams := false

		// Create a style for active parameters
		activeParamStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green color

		// Show users if any
		if users, ok := m.ParameterValues["users"]; ok && len(users) > 0 {
			hasLookupParams = true
			for _, user := range users {
				activeParamsContent.WriteString(activeParamStyle.Render(fmt.Sprintf("• user: %s\n", user)))
			}
		}

		// Show actions if any
		if actions, ok := m.ParameterValues["actions"]; ok && len(actions) > 0 {
			hasLookupParams = true
			// Group actions on one line to save space
			activeParamsContent.WriteString(activeParamStyle.Render("• actions: "))
			for i, action := range actions {
				if i > 0 {
					activeParamsContent.WriteString(", ")
				}
				activeParamsContent.WriteString(activeParamStyle.Render(action))
			}
			activeParamsContent.WriteString("\n")
		}

		// Show radius if any
		if radius, ok := m.MacroParameters["radius"]; ok && radius != "" {
			hasLookupParams = true
			activeParamsContent.WriteString(activeParamStyle.Render(fmt.Sprintf("• radius: %s\n", radius)))
		}

		// Show time if any
		if timeParam, ok := m.MacroParameters["time"]; ok && timeParam != "" {
			hasLookupParams = true
			activeParamsContent.WriteString(activeParamStyle.Render(fmt.Sprintf("• time: %s\n", timeParam)))
		}

		if !hasLookupParams {
			activeParamsContent.WriteString(m.SubtleStyle.Render("None\n"))
		}

		modalContent.WriteString(activeParamsStyle.Render(activeParamsContent.String()))
	}

	// Show help text
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	if selectedMacro.Name == "CoreProtect Pager" {
		modalContent.WriteString("\n" + helpStyle.Render("TAB: Next field • SHIFT+TAB: Previous field • CTRL+E: Add Parameter • ENTER: Confirm • ESC: Cancel"))
	} else {
		modalContent.WriteString("\n" + helpStyle.Render("TAB: Next field • SHIFT+TAB: Previous field • ENTER: Confirm • ESC: Cancel"))
	}

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
