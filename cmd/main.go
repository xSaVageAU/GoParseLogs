package main

import (
	"fmt"
	"os"

	"goparselogs/internal/macros"
	"goparselogs/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Register macros directly
	macros.RegisterHelloWorldMacro()
	macros.RegisterCoreProtectPagerMacro()

	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error running TUI: %v\n", err)
		os.Exit(1)
	}
}
