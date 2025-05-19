package fileops

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanLogFiles returns a list of .log and .log.gz files from the logs directory
func ScanLogFiles() ([]string, error) {
	// Get the logs directory relative to the current working directory
	logsDir := "logs"

	// Create logs directory if it doesn't exist
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		if err := os.Mkdir(logsDir, 0755); err != nil {
			return nil, err
		}
	}

	var files []string
	err := filepath.Walk(logsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			// Check for .log files and .gz files that end with .log.gz
			if ext == ".log" || (ext == ".gz" && strings.HasSuffix(strings.ToLower(path), ".log.gz")) {
				// Convert path separators to forward slashes for consistency
				normalizedPath := filepath.ToSlash(path)
				files = append(files, normalizedPath)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
