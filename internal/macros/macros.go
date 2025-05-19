package macros

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
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

// RegisterHelloWorldMacro registers the HelloWorld macro
func RegisterHelloWorldMacro() {
	RegisterMacro(Macro{
		Name:        "Type 'Hello World'",
		Description: "A macro that types 'Hello World' after a countdown.",
		Parameters:  []MacroParameter{}, // No parameters for this simple macro
		Action: func(params map[string]string) error {
			// Implement the HelloWorld macro directly here to avoid import cycles
			// This is a temporary solution until we restructure the code
			robotgo.TypeStr("hello world!")
			return nil
		},
	})
}

// RegisterCoreProtectPagerMacro registers the CoreProtect Pager macro
func RegisterCoreProtectPagerMacro() {
	RegisterMacro(Macro{
		Name:        "CoreProtect Pager",
		Description: "Runs /co lookup with parameters and then pages through results.",
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
			{
				Name:         "users",
				Description:  "Comma-separated list of usernames to lookup (e.g., 'user1,user2')",
				DefaultValue: "",
			},
			{
				Name:         "actions",
				Description:  "Comma-separated list of actions to lookup (e.g., 'block,chat')",
				DefaultValue: "",
			},
			{
				Name:         "radius",
				Description:  "Radius for the lookup (optional)",
				DefaultValue: "",
			},
			{
				Name:         "time",
				Description:  "Time parameter for lookup (e.g., '1d' for 1 day, '12h' for 12 hours)",
				DefaultValue: "",
			},
			{
				Name:         "skipLookup",
				Description:  "Set to 'true' to skip the lookup command and only run page commands",
				DefaultValue: "false",
			},
		},
		Action: func(params map[string]string) error {
			// Get required parameters
			startPageStr, okStart := params["startPage"]
			endPageStr, okEnd := params["endPage"]

			// Check if required parameters are present
			if !okStart || !okEnd {
				return fmt.Errorf("startPage and endPage parameters are required")
			}

			// Parse startPage parameter
			// Trim null bytes in case they were included during input
			startPageStr = strings.TrimRight(startPageStr, "\x00")
			startPage, err := strconv.Atoi(startPageStr)
			if err != nil {
				return fmt.Errorf("invalid startPage: %w", err)
			}

			// Parse endPage parameter
			// Trim null bytes in case they were included during input
			endPageStr = strings.TrimRight(endPageStr, "\x00")
			endPage, err := strconv.Atoi(endPageStr)
			if err != nil {
				return fmt.Errorf("invalid endPage: %w", err)
			}

			// Validate page range
			if startPage <= 0 || endPage < startPage {
				return fmt.Errorf("invalid page range: startPage must be > 0 and endPage >= startPage")
			}

			// Get optional delay parameter (default: 500ms)
			delayMs := 500 // Default delay
			if delayMsStr, okDelay := params["delayMs"]; okDelay {
				parsedDelay, err := strconv.Atoi(delayMsStr)
				if err == nil && parsedDelay > 0 {
					delayMs = parsedDelay
				}
			}

			// Check if we should execute the lookup command
			// We execute lookup if any of these parameters are present: users, actions, radius, time
			shouldExecuteLookup := false
			lookupCmd := "/co lookup"

			// Create a separate map for lookup parameters to ensure no interference
			lookupParams := make(map[string]string)

			// Explicitly get and add lookup parameters if specified
			if usersStr, ok := params["users"]; ok && usersStr != "" {
				shouldExecuteLookup = true
				lookupParams["users"] = usersStr
			}

			if actionsStr, ok := params["actions"]; ok && actionsStr != "" {
				shouldExecuteLookup = true
				lookupParams["actions"] = actionsStr
			}

			if radiusStr, ok := params["radius"]; ok && radiusStr != "" {
				shouldExecuteLookup = true
				lookupParams["radius"] = radiusStr
			}

			if timeStr, ok := params["time"]; ok && timeStr != "" {
				shouldExecuteLookup = true
				lookupParams["time"] = timeStr
			}

			// Build the lookup command string from the collected lookup parameters
			if usersStr, ok := lookupParams["users"]; ok {
				users := strings.Split(usersStr, ",")
				for _, user := range users {
					user = strings.TrimSpace(user)
					if user != "" {
						lookupCmd += fmt.Sprintf(" user:%s", user)
					}
				}
			}

			if actionsStr, ok := lookupParams["actions"]; ok {
				actions := strings.Split(actionsStr, ",")
				for _, action := range actions {
					action = strings.TrimSpace(action)
					if action != "" {
						lookupCmd += fmt.Sprintf(" action:%s", action)
					}
				}
			}

			if radiusStr, ok := lookupParams["radius"]; ok {
				lookupCmd += fmt.Sprintf(" r:%s", radiusStr)
			}

			if timeStr, ok := lookupParams["time"]; ok {
				lookupCmd += fmt.Sprintf(" t:%s", timeStr)
			}

			// Execute the lookup command if needed
			if shouldExecuteLookup {
				fmt.Printf("Executing lookup command: %s\n", lookupCmd)

				// Type the command
				robotgo.TypeStr(lookupCmd)

				// Press Enter to execute
				robotgo.KeyTap("enter")

				// Wait longer after the lookup command before paging
				time.Sleep(time.Duration(delayMs*2) * time.Millisecond)

				// Press 't' to open chat again for the first page command
				fmt.Println("Opening chat window for page commands...")
				robotgo.KeyTap("t")

				// Wait a moment after opening chat
				time.Sleep(100 * time.Millisecond)
			}

			// Execute the page commands
			fmt.Printf("Running CoreProtect pager from page %d to %d with %dms delay\n",
				startPage, endPage, delayMs)

			for i := startPage; i <= endPage; i++ {
				command := fmt.Sprintf("/co page %d", i)
				fmt.Printf("Executing: %s\n", command)

				// Type the command
				robotgo.TypeStr(command)

				// Press Enter to execute
				robotgo.KeyTap("enter")

				// Wait for a moment before opening chat again
				time.Sleep(100 * time.Millisecond)

				// Don't press 't' after the last command
				if i < endPage {
					// Press 't' to open chat again for the next command
					fmt.Println("Opening chat window for next command...")
					robotgo.KeyTap("t")

					// Wait a moment after opening chat
					time.Sleep(100 * time.Millisecond)
				}

				// Wait for the specified delay before the next command
				time.Sleep(time.Duration(delayMs) * time.Millisecond)
			}

			fmt.Println("CoreProtect pager completed successfully")
			return nil
		},
	})
}
