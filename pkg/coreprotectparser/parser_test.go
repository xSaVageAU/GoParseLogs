package coreprotectparser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseLogContent_CoreProtectEntries(t *testing.T) {
	logContent := `
[14:37:37] [Render thread/INFO]: [System] [CHAT] ----- CoreProtect |  Lookup Results -----
[14:37:37] [Render thread/INFO]: [System] [CHAT] 14.20/h ago §f- queercookie: §fcan I see?
[14:37:37] [Render thread/INFO]: [System] [CHAT] 10.5/h ago §f- Th4tGuy1: §fi dontr need it all
[14:37:37] [Render thread/INFO]: [System] [CHAT] CoreProtect - Lookup searching. Please wait...
[10:00:00] [Render thread/INFO]: [System] [CHAT] 1.0/h ago §f- anotheruser: §fsome message
[14:37:37] [Render thread/INFO]: [System] [CHAT] -----
[14:37:37] [Render thread/INFO]: [System] [CHAT] §f◀ Page §f101/6378 ▶ §7(§f1 §7... §f99 §7| §f100 §7| §f§n101§r §7| §f102 §7| §f103 §7... §f6378§7)
This is a regular log line, not CoreProtect.
[15:00:00] [Render thread/INFO]: [System] [CHAT] 0.5/h ago §f- user3: §fanother one
`
	expectedTime1, _ := time.Parse("15:04:05", "14:37:37")
	expectedTime2, _ := time.Parse("15:04:05", "10:00:00")
	expectedTime3, _ := time.Parse("15:04:05", "15:00:00")

	// Expected order is now oldest first (descending HoursAgo)
	// OriginalAcceptedIndex is not explicitly checked here as the primary sort key (HoursAgo) is distinct for these entries.
	expectedEntries := []CoreProtectLogEntry{
		{HoursAgo: 14.20, Username: "queercookie", Message: "can I see?", RawLine: "[14:37:37] [Render thread/INFO]: [System] [CHAT] 14.20/h ago §f- queercookie: §fcan I see?", Time: expectedTime1, OriginalAcceptedIndex: 0},
		{HoursAgo: 10.5, Username: "Th4tGuy1", Message: "i dontr need it all", RawLine: "[14:37:37] [Render thread/INFO]: [System] [CHAT] 10.5/h ago §f- Th4tGuy1: §fi dontr need it all", Time: expectedTime1, OriginalAcceptedIndex: 1},
		{HoursAgo: 1.0, Username: "anotheruser", Message: "some message", RawLine: "[10:00:00] [Render thread/INFO]: [System] [CHAT] 1.0/h ago §f- anotheruser: §fsome message", Time: expectedTime2, OriginalAcceptedIndex: 2},
		{HoursAgo: 0.5, Username: "user3", Message: "another one", RawLine: "[15:00:00] [Render thread/INFO]: [System] [CHAT] 0.5/h ago §f- user3: §fanother one", Time: expectedTime3, OriginalAcceptedIndex: 3},
	}

	parsedLog, err := ParseLogContent(logContent)
	assert.NoError(t, err)
	assert.NotNil(t, parsedLog)
	assert.Equal(t, len(expectedEntries), len(parsedLog.Entries), "Number of parsed entries does not match expected")

	for i, expected := range expectedEntries {
		actual := parsedLog.Entries[i]
		assert.Equal(t, expected.HoursAgo, actual.HoursAgo, "HoursAgo mismatch for entry %d", i)
		assert.Equal(t, expected.Username, actual.Username, "Username mismatch for entry %d", i)
		assert.Equal(t, expected.Message, actual.Message, "Message mismatch for entry %d", i)
		assert.Equal(t, expected.RawLine, actual.RawLine, "RawLine mismatch for entry %d", i)
		assert.True(t, expected.Time.Equal(actual.Time), "Time mismatch for entry %d. Expected %v, got %v", i, expected.Time, actual.Time)
		assert.Equal(t, expected.OriginalAcceptedIndex, actual.OriginalAcceptedIndex, "OriginalAcceptedIndex mismatch for entry %d", i)
	}
}

func TestParseLogContent_TiedHoursAgo(t *testing.T) {
	logContent := `
[10:00:00] [Render thread/INFO]: [System] [CHAT] 1.0/h ago §f- userA: §fmessage A (parsed first, younger in CP output for this tie)
[10:00:01] [Render thread/INFO]: [System] [CHAT] 0.5/h ago §f- userC: §fmessage C (youngest overall)
[10:00:02] [Render thread/INFO]: [System] [CHAT] 1.0/h ago §f- userB: §fmessage B (parsed second, older in CP output for this tie)
`
	timeA, _ := time.Parse("15:04:05", "10:00:00")
	timeC, _ := time.Parse("15:04:05", "10:00:01")
	timeB, _ := time.Parse("15:04:05", "10:00:02")

	// Expected:
	// 1. userB (1.0/h ago, OriginalAcceptedIndex 2 - parsed later, so older in the 1.0h tie)
	// 2. userA (1.0/h ago, OriginalAcceptedIndex 0 - parsed earlier, so younger in the 1.0h tie)
	// 3. userC (0.5/h ago, OriginalAcceptedIndex 1)
	expectedEntries := []CoreProtectLogEntry{
		{HoursAgo: 1.0, Username: "userB", Message: "message B (parsed second, older in CP output for this tie)", RawLine: "[10:00:02] [Render thread/INFO]: [System] [CHAT] 1.0/h ago §f- userB: §fmessage B (parsed second, older in CP output for this tie)", Time: timeB, OriginalAcceptedIndex: 2},
		{HoursAgo: 1.0, Username: "userA", Message: "message A (parsed first, younger in CP output for this tie)", RawLine: "[10:00:00] [Render thread/INFO]: [System] [CHAT] 1.0/h ago §f- userA: §fmessage A (parsed first, younger in CP output for this tie)", Time: timeA, OriginalAcceptedIndex: 0},
		{HoursAgo: 0.5, Username: "userC", Message: "message C (youngest overall)", RawLine: "[10:00:01] [Render thread/INFO]: [System] [CHAT] 0.5/h ago §f- userC: §fmessage C (youngest overall)", Time: timeC, OriginalAcceptedIndex: 1},
	}

	parsedLog, err := ParseLogContent(logContent)
	assert.NoError(t, err)
	assert.NotNil(t, parsedLog)
	assert.Equal(t, len(expectedEntries), len(parsedLog.Entries), "Number of parsed entries does not match expected for tied HoursAgo")

	for i, expected := range expectedEntries {
		actual := parsedLog.Entries[i]
		assert.Equal(t, expected.HoursAgo, actual.HoursAgo, "Tied HoursAgo: HoursAgo mismatch for entry %d", i)
		assert.Equal(t, expected.Username, actual.Username, "Tied HoursAgo: Username mismatch for entry %d", i)
		assert.Equal(t, expected.Message, actual.Message, "Tied HoursAgo: Message mismatch for entry %d", i)
		// OriginalAcceptedIndex is key for this test's tie-breaking logic.
		assert.Equal(t, expected.OriginalAcceptedIndex, actual.OriginalAcceptedIndex, "Tied HoursAgo: OriginalAcceptedIndex mismatch for entry %d", i)
	}
}

func TestParseLogContent_EmptyLog(t *testing.T) {
	logContent := ""
	parsedLog, err := ParseLogContent(logContent)
	assert.NoError(t, err)
	assert.NotNil(t, parsedLog)
	assert.Empty(t, parsedLog.Entries)
}

func TestParseLogContent_NoCoreProtectEntries(t *testing.T) {
	logContent := `
[14:37:37] [Render thread/INFO]: This is a standard log line.
[14:37:38] [Render thread/INFO]: Another standard log line.
`
	parsedLog, err := ParseLogContent(logContent)
	assert.NoError(t, err)
	assert.NotNil(t, parsedLog)
	assert.Empty(t, parsedLog.Entries)
}

func TestParseLogContent_MalformedHours(t *testing.T) {
	logContent := `[14:37:37] [Render thread/INFO]: [System] [CHAT] nothours/h ago §f- queercookie: §fcan I see?`
	parsedLog, err := ParseLogContent(logContent)
	assert.NoError(t, err)
	assert.NotNil(t, parsedLog)
	assert.Empty(t, parsedLog.Entries, "Should not parse entries with malformed hours")
}
