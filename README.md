# GoParseLogs

A terminal-based log viewer for parsing and filtering log files. Supports both plain text logs and gzipped log files.

![Log Viewer Demo](demo.gif)

## Features

- 📂 Automatically detects `.log` and `.log.gz` files in the logs directory
- 🔍 Real-time filtering of log entries
- 💾 Save filtered results to output files
- 📦 Handles gzipped log files seamlessly
- 🔄 Auto-refreshes when new log files are added
- ⚡ CoreProtect log parsing support

## Requirements

### Pre-built Package
Download the latest release from the [releases page](releases) for your platform.

### Building from Source
- Go 1.20 or higher
- Run `go build ./cmd/main.go` to build
- Or use `go run cmd/main.go` to run directly

## Usage

1. Place your log files in the `logs` directory
2. Run the application:
   - If using pre-built: `./goparselogs` (or `goparselogs.exe` on Windows)
   - If running from source: `go run cmd/main.go`
3. Navigate files with UP/DOWN or J/K
4. Press TAB to add filters
5. Press E to export filtered results
6. Press Q or Ctrl+C to quit

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