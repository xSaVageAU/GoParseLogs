package macros

// Macro represents a defined macro.
type Macro struct {
	Name        string
	Description string
	// Later, this might include the actual sequence of actions.
}

// GetAvailableMacros returns a list of currently defined macros.
// For now, these are placeholders.
func GetAvailableMacros() []Macro {
	return []Macro{
		{Name: "Example Macro 1", Description: "This is the first example macro."},
		{Name: "Example Macro 2", Description: "Another macro for demonstration."},
		{Name: "Type 'Hello World'", Description: "A macro that will type 'Hello World'."},
	}
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
