package main

import (
	"fmt"
	"os"

	"goparselogs/internal/macros"
	"goparselogs/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("GoParseLogs - Starting application...")
	fmt.Println("Note: This application uses keyboard automation which may trigger antivirus warnings.")
	fmt.Println("If you're experiencing issues, consider adding an exception in your antivirus software.")

	// Register macros directly
	macros.RegisterHelloWorldMacro()
	macros.RegisterCoreProtectPagerMacro()

	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error running TUI: %v\n", err)
		os.Exit(1)
	}
}
