package macros

import (
	"fmt"
	"goparselogs/internal/macros/scripts"
	"os"      // For reading directory contents
	"strings" // For string manipulation
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

// GetAvailableMacros scans the scripts directory and returns a list of recognized macros.
func GetAvailableMacros() []Macro {
	var availableMacros []Macro
	scriptsDir := "internal/macros/scripts" // Relative to project root

	files, err := os.ReadDir(scriptsDir)
	if err != nil {
		// If the directory can't be read, return an empty list or handle error
		// For now, print an error to stderr and return empty, so the app doesn't crash.
		fmt.Fprintf(os.Stderr, "Error reading scripts directory '%s': %v\n", scriptsDir, err)
		return availableMacros
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			scriptName := strings.TrimSuffix(file.Name(), ".go")
			// This mapping needs to be maintained as new scripts are added.
			switch scriptName {
			case "helloworld":
				availableMacros = append(availableMacros, Macro{
					Name:        "Type 'Hello World'",
					Description: "A macro that types 'Hello World' after a countdown.",
					Parameters:  []MacroParameter{}, // No parameters for this simple macro
					Action: func(params map[string]string) error {
						return scripts.RunHelloWorld()
					},
				})
			case "coreprotect_pager":
				availableMacros = append(availableMacros, Macro{
					Name:        "CoreProtect Pager",
					Description: "Runs /co page X commands from a start page to an end page.",
					Parameters: []MacroParameter{
						{
							Name:         "startPage",
							Description:  "Starting page number",
							DefaultValue: "1",
						},
						{
							Name:         "endPage",
							Description:  "Ending page number",
							DefaultValue: "5",
						},
						{
							Name:         "delayMs",
							Description:  "Delay in milliseconds between commands (optional)",
							DefaultValue: "500",
						},
					},
					Action: scripts.RunCoreProtectPager,
				})
				// Add other cases here for new scripts:
				// case "another_script_filename_without_ext":
				//	 availableMacros = append(availableMacros, Macro{
				//		 Name:        "User Friendly Name for Another Script",
				//		 Description: "What this other script does.",
				//		 Action:      scripts.RunAnotherScriptFunction, // Assuming scripts.RunAnotherScriptFunction exists
				//	 })
			}
		}
	}
	return availableMacros
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
