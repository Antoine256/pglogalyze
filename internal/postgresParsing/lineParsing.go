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

type LType string

const (
	QUERY      LType = "QUERY"
	CONNECTION LType = "CONNECTION"
	DURATION   LType = "DURATION"
	CHECKPOINT LType = "CHECKPOINT"
	STARTUP    LType = "STARTUP"
	SHUTDOWN   LType = "SHUTDOWN"
	ALL        LType = ""
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

func IsAValidType(s string) bool {
	s = strings.ToUpper(s)
	switch s {
	case "QUERY", "CONNECTION", "DURATION", "CHECKPOINT", "STARTUP", "SHUTDOWN":
		return true
	default:
		return false
	}
}

type ParsedLineType struct {
	date          string
	hour          string
	time          time.Time
	timeZone      string
	pid           string
	bddInfo       string
	severity      Severity
	severityColor internal.CMDCOLOR
	logtype       LType
	logMessage    string
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
	lineDate := splitLine[0]
	lineSeverity := Severity(strings.Replace(splitLine[4], ":", "", 1))
	lineMessage := strings.Join(splitLine[5:], " ")
	lineBDDInfo := ""

	//si c'est une ligne appli elle commence par [\]
	if line[:3] == "[\\]" {
		lineDate = splitLine[0][3:]
		lineBDDInfo = splitLine[4]
		lineSeverity = Severity(strings.Replace(splitLine[5], ":", "", 1))
		lineMessage = strings.Join(splitLine[6:], " ")
	}

	lineSeverityColor := getLineSeverityColor(lineSeverity)
	linelogtype := getlineLogType(lineSeverity, strings.TrimLeft(lineMessage, " "))

	parsedLine := ParsedLineType{
		date:          lineDate,
		hour:          splitLine[1],
		timeZone:      splitLine[2],
		pid:           splitLine[3],
		bddInfo:       lineBDDInfo,
		severity:      lineSeverity,
		severityColor: lineSeverityColor,
		logtype:       linelogtype,
		logMessage:    lineMessage,
	}

	parsedLine.time = internal.StringToTime(lineDate, splitLine[1])

	return parsedLine
}

func getLineSeverityColor(lineSeverity Severity) internal.CMDCOLOR {
	switch lineSeverity {
	case DEBUG1, DEBUG2, DEBUG3, DEBUG4, DEBUG5, WARNING:
		return internal.Yellow
	case INFO, LOG, NOTICE, STATEMENT, DETAIL, NONE:
		return internal.Blue
	case ERROR, FATAL, PANIC:
		return internal.Red
	default:
		return internal.Reset
	}
}

func getlineLogType(severity Severity, message string) LType {

	// QUERY
	//If the severity is STATEMENT or the first word of the message is statement
	firstWord := getFirstWordOfString(message)

	if severity == Severity("STATEMENT") || firstWord == "statement" {
		return QUERY
	}

	if strings.Contains(message, "starting") || strings.Contains(message, "listening on") || strings.Contains(message, "was shut down") || strings.Contains(message, "is ready") {
		return STARTUP
	}
	if strings.Contains(strings.Trim(message, " "), "shutdown") || strings.Contains(message, "aborting") || strings.Contains(message, "exited with exit code") || strings.Contains(message, "shutting down") {
		return SHUTDOWN
	}

	switch firstWord {
	case "duration":
		return DURATION
	case "connection", "disconnection":
		return CONNECTION
	case "checkpoint":
		return CHECKPOINT
	default:
		return ALL
	}
}

func getFirstWordOfString(chaine string) string {
	str := ""
	for i := 0; i < len(chaine); i++ {
		if strings.Contains(" :/;'&-()_?,.§!", string(chaine[i])) {
			break
		}
		str = str + string(chaine[i])
	}

	return str
}
