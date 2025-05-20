package coreprotectparser

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// CoreProtectLogEntry represents a parsed CoreProtect log entry.
type CoreProtectLogEntry struct {
	HoursAgo              float64 // Total time ago in hours (converted from days if needed)
	DaysAgo               float64 // Original days value if present, 0 otherwise
	IsInDays              bool    // Whether the original time was in days
	Username              string
	Message               string
	RawLine               string
	Time                  time.Time // Parsed from the initial log timestamp
	OriginalAcceptedIndex int       // The order in which this entry was accepted during parsing
}

// ParsedLog represents the overall parsed log with CoreProtect entries.
type ParsedLog struct {
	Entries []CoreProtectLogEntry
}

var (
	// Regex to capture the hours ago, username, and message from CoreProtect chat lines
	// Example: 14.20/h ago §f- queercookie: §fcan I see?
	coreProtectHoursChatRegex = regexp.MustCompile(`\[(\d{2}:\d{2}:\d{2})\] \[Render thread/INFO\]: \[System\] \[CHAT\] (\d+\.\d+)/h ago §f- ([^:]+): §f(.*)`)

	// Regex to capture the days ago, username, and message from CoreProtect chat lines
	// Example: 1023.18/d ago §f- xSaVage_: §f0.05
	coreProtectDaysChatRegex = regexp.MustCompile(`\[(\d{2}:\d{2}:\d{2})\] \[Render thread/INFO\]: \[System\] \[CHAT\] (\d+\.\d+)/d ago §f- ([^:]+): §f(.*)`)

	// Regex to identify lines that are part of CoreProtect lookup but not actual messages
	coreProtectMetaRegex = regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\] \[Render thread/INFO\]: \[System\] \[CHAT\] (----- CoreProtect \| Lookup Results -----|CoreProtect - Lookup searching\. Please wait\.\.\.|§f◀ Page §f\d+/\d+ ▶)`)
)

// ParseLogContent parses the raw log content string and extracts CoreProtect entries.
func ParseLogContent(logContent string) (*ParsedLog, error) {
	lines := strings.Split(logContent, "\n")
	parsedLog := &ParsedLog{}
	acceptedCounter := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Skip lines that are CoreProtect metadata but not chat messages
		if coreProtectMetaRegex.MatchString(line) {
			continue
		}

		// Try to match hours format first
		match := coreProtectHoursChatRegex.FindStringSubmatch(line)
		if len(match) == 5 {
			timestampStr := match[1]
			hoursAgoStr := match[2]
			username := strings.TrimSpace(match[3])
			message := strings.TrimSpace(match[4])

			hoursAgo, err := strconv.ParseFloat(hoursAgoStr, 64)
			if err != nil {
				// Skip lines where hours ago cannot be parsed
				continue
			}

			logTime, err := time.Parse("15:04:05", timestampStr)
			if err != nil {
				logTime = time.Time{}
			}

			parsedLog.Entries = append(parsedLog.Entries, CoreProtectLogEntry{
				HoursAgo:              hoursAgo,
				DaysAgo:               0,
				IsInDays:              false,
				Username:              username,
				Message:               message,
				RawLine:               line,
				Time:                  logTime,
				OriginalAcceptedIndex: acceptedCounter,
			})
			acceptedCounter++
			continue
		}

		// If not hours, try to match days format
		match = coreProtectDaysChatRegex.FindStringSubmatch(line)
		if len(match) == 5 {
			timestampStr := match[1]
			daysAgoStr := match[2]
			username := strings.TrimSpace(match[3])
			message := strings.TrimSpace(match[4])

			daysAgo, err := strconv.ParseFloat(daysAgoStr, 64)
			if err != nil {
				// Skip lines where days ago cannot be parsed
				continue
			}

			// Convert days to hours for consistent sorting
			hoursAgo := daysAgo * 24

			logTime, err := time.Parse("15:04:05", timestampStr)
			if err != nil {
				logTime = time.Time{}
			}

			parsedLog.Entries = append(parsedLog.Entries, CoreProtectLogEntry{
				HoursAgo:              hoursAgo,
				DaysAgo:               daysAgo,
				IsInDays:              true,
				Username:              username,
				Message:               message,
				RawLine:               line,
				Time:                  logTime,
				OriginalAcceptedIndex: acceptedCounter,
			})
			acceptedCounter++
		}
	}

	// Sort entries:
	// 1. By HoursAgo in descending order (oldest "HoursAgo" value first).
	// 2. If HoursAgo are equal, by OriginalAcceptedIndex in descending order
	//    (entry that appeared later in the log, thus "older" in a tied group, comes first).
	sort.Slice(parsedLog.Entries, func(i, j int) bool {
		entryI := parsedLog.Entries[i]
		entryJ := parsedLog.Entries[j]

		if entryI.HoursAgo != entryJ.HoursAgo {
			return entryI.HoursAgo > entryJ.HoursAgo
		}
		return entryI.OriginalAcceptedIndex > entryJ.OriginalAcceptedIndex
	})

	return parsedLog, nil
}