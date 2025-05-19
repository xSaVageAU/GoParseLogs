# GoParseLogs

A terminal-based log viewer for parsing and filtering Minecraft logs. Works with both client and server log files, supporting plain text and gzipped formats.

![GoParseLogs](https://github.com/user-attachments/assets/ca0fd566-364e-47fb-8335-fb104d556140)

## Features

- 📂 Automatically detects `.log` and `.log.gz` files in the logs directory
- 🔍 Real-time filtering of log entries
- 💾 Save filtered results to output files
- 📦 Handles gzipped log files seamlessly
- 🔄 Auto-refreshes when new log files are added
- ⚡ CoreProtect log parsing support

## Requirements

### Pre-built Package (Windows Only)
Download the latest Windows release from the [releases page](https://github.com/xSaVageAU/GoParseLogs/releases).

### Building from Source
- Go 1.20 or higher
- Windows OS
- Run `go build ./cmd/main.go` to build
- Or use `go run cmd/main.go` to run directly

## Usage

You can either:
1. Place the application in your Minecraft client or server folder (it will automatically read from the existing `logs` directory)
2. Or create a new `logs` directory next to the executable and place your log files there

Then:
1. Run the application:
   - If using pre-built: `goparselogs.exe`
   - If running from source: `go run cmd/main.go`
2. Navigate files with UP/DOWN or J/K
3. Press TAB to add filters
4. Press E to export filtered results
5. Press Q or Ctrl+C to quit

### Keyboard Shortcuts

- `↑/↓` or `j/k`: Navigate logs
- `Tab`: Toggle between files and filter input
- `Enter`: Select file / Apply filter
- `e`: Export filtered logs
- `q` or `Ctrl+C`: Quit

## AI Disclaimer

This project was developed with the assistance of AI (Claude). The entire codebase, including structure and implementation, was generated through AI-driven development while maintaining high code quality and following Go best practices.

---
Made with ❤️ and 🤖 (Claude)
