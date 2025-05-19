package scripts

import "github.com/go-vgo/robotgo"

// RunHelloWorld types "hello world!" using robotgo.
func RunHelloWorld() error {
	robotgo.TypeStr("hello world!")
	return nil
}
