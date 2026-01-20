package postgresparsing

import (
	"pglogalyze/internal"
	"strings"
	"time"
)

type Severity string

const (
	DEBUG1    Severity = "DEBUG1"
	DEBUG2    Severity = "DEBUG2"
	DEBUG3    Severity = "DEBUG3"
	DEBUG4    Severity = "DEBUG4"
	DEBUG5    Severity = "DEBUG5"
	INFO      Severity = "INFO"
	NOTICE    Severity = "NOTICE"
	WARNING   Severity = "WARNING"
	ERROR     Severity = "ERROR"
	LOG       Severity = "LOG"
	FATAL     Severity = "FATAL"
	PANIC     Severity = "PANIC"
	STATEMENT Severity = "STATEMENT"
	DETAIL    Severity = "DETAIL"
	NONE      Severity = ""
)

func IsAValidSeverity(s string) bool {
	s = strings.ToUpper(s)
	switch s {
	case "DEBUG1", "DEBUG2", "DEBUG3", "DEBUG4", "DEBUG5", "INFO", "NOTICE", "WARNING", "ERROR", "LOG", "FATAL", "PANIC", "STATEMENT", "DETAIL":
		return true
	default:
		return false
	}
}

type ParsedLineType struct {
	date       string
	hour       string
	time       time.Time
	timeZone   string
	pid        string
	severity   Severity
	logMessage string
}

func parseLines(lines []string) []ParsedLineType {

	parsedLogs := []ParsedLineType{}

	for _, line := range lines {
		parsedLine := parseLine(line)
		parsedLogs = append(parsedLogs, parsedLine)
	}

	return parsedLogs

}

func parseLine(line string) ParsedLineType {
	line = strings.TrimSpace(line)
	splitLine := strings.Split(line, " ")

	splitLine[4] = strings.Replace(splitLine[4], ":", "", 1)

	parsedLine := ParsedLineType{
		date:       splitLine[0],
		hour:       splitLine[1],
		timeZone:   splitLine[2],
		pid:        splitLine[3],
		severity:   Severity(splitLine[4]),
		logMessage: strings.Join(splitLine[5:], " "),
	}

	parsedLine.time = internal.StringToTime(splitLine[0], splitLine[1])

	return parsedLine
}
