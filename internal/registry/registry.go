package registry

import (
	"fmt"
)

// MacroParameter represents a parameter for a macro.
type MacroParameter struct {
	Name         string
	Description  string
	DefaultValue string // Optional default value
}

// Macro represents a defined macro.
type Macro struct {
	Name        string
	Description string
	Parameters  []MacroParameter                     // Parameters that the macro accepts
	Action      func(params map[string]string) error // Function to execute for the macro with parameters
}

// Registry of all available macros
var registeredMacros []Macro

// RegisterMacro adds a macro to the registry
func RegisterMacro(macro Macro) {
	fmt.Printf("Registering macro: %s\n", macro.Name)
	registeredMacros = append(registeredMacros, macro)
}

// GetAvailableMacros returns the list of registered macros
func GetAvailableMacros() []Macro {
	return registeredMacros
}

// GetMacroNames returns just the names of the available macros.
func GetMacroNames() []string {
	availableMacros := GetAvailableMacros()
	names := make([]string, len(availableMacros))
	for i, m := range availableMacros {
		names[i] = m.Name
	}
	return names
}

// ExecuteMacro finds a macro by name and executes its action with the provided parameters.
func ExecuteMacro(macroName string, params map[string]string) error {
	fmt.Printf("Executing macro: %s with params: %v\n", macroName, params)
	var selectedMacro *Macro
	for _, m := range GetAvailableMacros() {
		if m.Name == macroName {
			// Create a new variable for the loop to avoid capturing the loop variable in the closure.
			tempMacro := m
			selectedMacro = &tempMacro
			break
		}
	}

	if selectedMacro == nil {
		return fmt.Errorf("macro '%s' not found", macroName)
	}

	if selectedMacro.Action == nil {
		return fmt.Errorf("macro '%s' has no action defined", macroName)
	}

	// Apply default values for any parameters that weren't provided
	for _, param := range selectedMacro.Parameters {
		if _, exists := params[param.Name]; !exists && param.DefaultValue != "" {
			params[param.Name] = param.DefaultValue
		}
	}

	// Execute the macro with the parameters
	return selectedMacro.Action(params)
}

// InitRegistry initializes the macro registry
func InitRegistry() {
	// Clear any existing registrations
	registeredMacros = nil

	// Debug: Print the number of registered macros
	fmt.Printf("Initialized macro registry with %d macros\n", len(registeredMacros))
}
