# Macro Scripts

This directory contains macro scripts for the GoParseLogs application. Each script is a Go file that defines a macro that can be executed from the application.

## How Macros Are Registered

Macros are registered directly in the application. The registration process works as follows:

1. Create a registration function in the `internal/macros/macros.go` file
2. Call this registration function from `main.go`
3. The macro registry collects all registered macros and makes them available to the UI

## Creating a New Macro

To create a new macro:

1. Add a new registration function in `internal/macros/macros.go`:
   ```go
   func RegisterNewMacro() {
       RegisterMacro(Macro{
           Name:        "New Macro Name",
           Description: "Description of what this macro does",
           Parameters:  []MacroParameter{},
           Action: func(params map[string]string) error {
               // Implement the macro functionality here
               return nil
           },
       })
   }
   ```

2. Call this registration function in `main.go`:
   ```go
   func main() {
       // Register macros directly
       macros.RegisterHelloWorldMacro()
       macros.RegisterCoreProtectPagerMacro()
       macros.RegisterNewMacro() // Add your new macro here
       
       p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
       // ...
   }
   ```

## Macro Structure

Each macro consists of:

- **Name**: A user-friendly name shown in the UI
- **Description**: A detailed description of what the macro does
- **Parameters**: A list of parameters that the macro accepts
- **Action**: A function that implements the macro logic

## Example

```go
// In internal/macros/macros.go
func RegisterMyCustomMacro() {
	RegisterMacro(Macro{
		Name:        "My Custom Macro",
		Description: "Description of what this macro does",
		Parameters: []MacroParameter{
			{
				Name:         "param1",
				Description:  "Description of parameter 1",
				DefaultValue: "default",
			},
		},
		Action: func(params map[string]string) error {
			// Get parameters
			param1, ok := params["param1"]
			if !ok {
				return fmt.Errorf("param1 is required")
			}
			
			// Your implementation here
			fmt.Printf("Running macro with param1: %s\n", param1)
			return nil
		},
	})
}

// In cmd/main.go
func main() {
	// Register macros
	macros.RegisterHelloWorldMacro()
	macros.RegisterCoreProtectPagerMacro()
	macros.RegisterMyCustomMacro()
	
	// ...
}
```

## Best Practices

1. Give your macro a clear, descriptive name
2. Provide a detailed description of what the macro does
3. Document each parameter with a clear description
4. Provide sensible default values for parameters when appropriate
5. Handle errors gracefully and return meaningful error messages
6. Follow Go best practices for code organization and documentation