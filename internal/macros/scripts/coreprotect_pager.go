package scripts

import (
	"fmt"
	"goparselogs/internal/macros"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"
)

// RunCoreProtectPager executes a series of "/co page X" commands from startPage to endPage
// with a configurable delay between commands.
func RunCoreProtectPager(params map[string]string) error {
	// Get required parameters
	startPageStr, okStart := params["startPage"]
	endPageStr, okEnd := params["endPage"]

	// Check if required parameters are present
	if !okStart || !okEnd {
		return fmt.Errorf("startPage and endPage parameters are required")
	}

	// Parse startPage parameter
	startPage, err := strconv.Atoi(startPageStr)
	if err != nil {
		return fmt.Errorf("invalid startPage: %w", err)
	}

	// Parse endPage parameter
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

	// Execute the commands
	fmt.Printf("Running CoreProtect pager from page %d to %d with %dms delay\n",
		startPage, endPage, delayMs)

	for i := startPage; i <= endPage; i++ {
		command := fmt.Sprintf("/co page %d", i)
		fmt.Printf("Executing: %s\n", command)

		// Type the command
		robotgo.TypeStr(command)

		// Press Enter to execute
		robotgo.KeyTap("enter")

		// Wait for the specified delay before the next command
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}

	fmt.Println("CoreProtect pager completed successfully")
	return nil
}

// init registers this macro with the macro registry
func init() {
	macros.RegisterMacro(macros.Macro{
		Name:        "CoreProtect Pager",
		Description: "Runs /co page X commands from a start page to an end page.",
		Parameters: []macros.MacroParameter{
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
		Action: RunCoreProtectPager,
	})
}
