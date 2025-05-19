package macros

import (
	"fmt"
	"goparselogs/internal/macros/scripts"
	"os"      // For reading directory contents
	"strings" // For string manipulation
	"time"
)

// Macro represents a defined macro.
type Macro struct {
	Name        string
	Description string
	Action      func() error // Function to execute for the macro
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
					Action:      scripts.RunHelloWorld,
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

// ExecuteMacro finds a macro by name and executes its action with a countdown.
func ExecuteMacro(macroName string) error {
	fmt.Printf("Executing macro: %s\n", macroName)
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

	fmt.Println("Starting in...")
	for i := 10; i > 0; i-- {
		fmt.Printf("%d...\n", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Executing now!")

	return selectedMacro.Action()
}
