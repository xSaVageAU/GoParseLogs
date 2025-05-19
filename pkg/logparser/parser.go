package logparser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// LogEntry represents a single parsed log entry.
type LogEntry struct {
	Timestamp string
	Thread    string
	Level     string
	Message   string
}

// Parser is responsible for parsing log files.
type Parser struct {
	logRegex *regexp.Regexp
}

// NewParser creates a new instance of the log parser.
func NewParser() (*Parser, error) {
	// Regex to capture timestamp, thread/level, and message
	regex, err := regexp.Compile(`^\[(\d{2}:\d{2}:\d{2})\] \[([^/]+)/([^\]]+)\]: (.*)$`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}
	return &Parser{logRegex: regex}, nil
}

// ParseLine parses a single line and returns a LogEntry if it matches the log format
func (p *Parser) ParseLine(line string) (LogEntry, error) {
	matches := p.logRegex.FindStringSubmatch(line)
	if len(matches) != 5 {
		return LogEntry{}, fmt.Errorf("line does not match log format: %s", line)
	}

	return LogEntry{
		Timestamp: matches[1],
		Thread:    matches[2],
		Level:     matches[3],
		Message:   matches[4],
	}, nil
}

// ParseContent parses log content from a string, applying filters
func (p *Parser) ParseContent(content string, filters []string) ([]LogEntry, error) {
	var entries []LogEntry
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		entry, err := p.ParseLine(line)
		if err != nil {
			continue // Skip lines that don't match the log format
		}

		// Apply filters using OR logic
		if len(filters) > 0 {
			matchesAtLeastOneFilter := false
			for _, filter := range filters {
				filterLower := strings.ToLower(filter)
				if strings.Contains(strings.ToLower(entry.Message), filterLower) ||
					strings.Contains(strings.ToLower(entry.Thread), filterLower) ||
					strings.Contains(strings.ToLower(entry.Level), filterLower) ||
					strings.Contains(strings.ToLower(entry.Timestamp), filterLower) {
					matchesAtLeastOneFilter = true
					break // Found a match, no need to check other filters for this entry
				}
			}
			if !matchesAtLeastOneFilter {
				continue // Skip this entry if it doesn't match any of the filters
			}
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log content: %w", err)
	}

	return entries, nil
}

// ParseLog extracts information from a log file, applying filters.
func (p *Parser) ParseLog(logFilePath string, filters []string) ([]LogEntry, error) {
	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		entry, err := p.ParseLine(line)
		if err != nil {
			continue // Skip lines that don't match the log format
		}

		// Apply filters using OR logic
		if len(filters) > 0 {
			matchesAtLeastOneFilter := false
			for _, filter := range filters {
				filterLower := strings.ToLower(filter)
				if strings.Contains(strings.ToLower(entry.Message), filterLower) ||
					strings.Contains(strings.ToLower(entry.Thread), filterLower) ||
					strings.Contains(strings.ToLower(entry.Level), filterLower) ||
					strings.Contains(strings.ToLower(entry.Timestamp), filterLower) {
					matchesAtLeastOneFilter = true
					break // Found a match, no need to check other filters for this entry
				}
			}
			if !matchesAtLeastOneFilter {
				continue // Skip this entry if it doesn't match any of the filters
			}
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	return entries, nil
}
