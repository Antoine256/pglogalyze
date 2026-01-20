package postgresparsing

import (
	"fmt"
	"os"
	"pglogalyze/internal"
	"time"
)

type Options struct {
	LogFilePath string
	Level       Severity
	StartTime   *time.Time
	EndTime     *time.Time
}

func ReadLogFile(options Options) {
	osfile, err := os.Open(options.LogFilePath)

	if err != nil {
		internal.PrintError("internal\\postgresParsing\\readingfile.go", err.Error())
	}

	file, err := os.ReadFile(osfile.Name())

	if err != nil {
		internal.PrintError("internal\\postgresParsing\\readingfile.go", err.Error())
	}

	//get Lines of the file as []string
	lines := parseFile(file)

	//parse the Lines to ParsedLineType
	parsedLines := parseLines(lines)

	//Sort lines depending on options
	sortedParsedLines := sortLogsLines(parsedLines, options)

	for i := 0; i < len(sortedParsedLines); i++ {
		fmt.Println(sortedParsedLines[i].toString())
	}

}

func (l ParsedLineType) toString() string {
	return l.date + " " + l.hour + " " + l.timeZone + " " + l.pid + " " + string(l.severity) + ":" + l.logMessage
}

func sortLogsLines(lines []ParsedLineType, options Options) []ParsedLineType {

	// Start and End time parameters

	if options.StartTime != nil {
		//RECHERCHE DICHO ???
		i := 0
		for i < len(lines) {
			if lines[i].time.Compare(*options.StartTime) >= 0 {
				break
			} else {
				i++
			}
		}
		if i == len(lines) {
			internal.PrintInfo("No line found after the start time " + options.StartTime.Format(time.UnixDate))
			return []ParsedLineType{}
		}
		lines = lines[i:]
	}

	if options.EndTime != nil {
		i := len(lines) - 1
		for i > 0 {
			if lines[i].time.Compare(*options.EndTime) <= 0 {
				break
			} else {
				i--
			}
		}
		if i == 0 {
			internal.PrintInfo("No line found before the end time " + options.EndTime.Format(time.UnixDate))
			return []ParsedLineType{}
		}
		lines = lines[:i+1]
	}

	// Severity parameter

	if options.Level != NONE {
		severityLines := []ParsedLineType{}
		for i := range lines {
			if lines[i].severity == options.Level {
				severityLines = append(severityLines, lines[i])
			}
		}
		lines = severityLines
	}

	return lines
}
