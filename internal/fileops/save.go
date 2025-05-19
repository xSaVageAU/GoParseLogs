package fileops

import (
	"fmt"
	"os"
	"strings"

	"goparselogs/pkg/coreprotectparser"
	"goparselogs/pkg/logparser"
)

// SaveStandardLogsToFile writes the provided standard log entries to a file.
func SaveStandardLogsToFile(entries []logparser.LogEntry, filename string) error {
	if len(entries) == 0 {
		return fmt.Errorf("no entries to save")
	}

	var contentBuilder strings.Builder
	for _, entry := range entries {
		contentBuilder.WriteString(fmt.Sprintf("[%s] [%s/%s]: %s\n", entry.Timestamp, entry.Thread, entry.Level, entry.Message))
	}

	// Create the output directory if it doesn't exist
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	filePath := fmt.Sprintf("%s/%s", outputDir, filename)
	err := os.WriteFile(filePath, []byte(contentBuilder.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	return nil
}

// SaveCoreProtectLogsToFile writes the provided CoreProtect log entries to a file.
func SaveCoreProtectLogsToFile(entries []coreprotectparser.CoreProtectLogEntry, filename string) error {
	if len(entries) == 0 {
		return fmt.Errorf("no CoreProtect entries to save")
	}

	var contentBuilder strings.Builder
	for _, entry := range entries {
		contentBuilder.WriteString(entry.RawLine + "\n")
	}

	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	filePath := fmt.Sprintf("%s/%s", outputDir, filename)
	err := os.WriteFile(filePath, []byte(contentBuilder.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}
	return nil
}
