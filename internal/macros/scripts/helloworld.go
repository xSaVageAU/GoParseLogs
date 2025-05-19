package scripts

import (
	"goparselogs/internal/macros"

	"github.com/go-vgo/robotgo"
)

// RunHelloWorld types "hello world!" using robotgo.
func RunHelloWorld() error {
	robotgo.TypeStr("hello world!")
	return nil
}

// init registers this macro with the macro registry
func init() {
	macros.RegisterMacro(macros.Macro{
		Name:        "Type 'Hello World'",
		Description: "A macro that types 'Hello World' after a countdown.",
		Parameters:  []macros.MacroParameter{}, // No parameters for this simple macro
		Action: func(params map[string]string) error {
			return RunHelloWorld()
		},
	})
}
