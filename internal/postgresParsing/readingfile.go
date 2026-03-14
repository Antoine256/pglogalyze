package postgresparsing

import (
	"fmt"
	"os"
	"pglogalyze/internal"
	"strconv"
	"time"
)

type Options struct {
	LogFilePath string
	Level       Severity
	LogType     LType
	StartTime   *time.Time
	EndTime     *time.Time
	NBLines     int
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
	//A OPTIMISER > parse toutes les lignes avec toutes les infos (pas en fonction des paramètres)
	parsedLines := parseLines(lines)

	//Sort lines depending on options
	sortedParsedLines := sortLogsLines(parsedLines, options)

	fmt.Println("Lines found : " + internal.Yellow.String() + strconv.Itoa(len(sortedParsedLines)) + "/" + strconv.Itoa(options.NBLines) + internal.Reset.String())

	for i := 0; i < len(sortedParsedLines); i++ {
		//Permet d'aligner les lignes après le PID
		if len(sortedParsedLines[i].pid) < 7 {
			for j := 0; j < 7; j++ {
				if len(sortedParsedLines[i].pid) == 7 {
					break
				}
				sortedParsedLines[i].pid = sortedParsedLines[i].pid + " "
			}
		}
		if sortedParsedLines[i].bddInfo != "" {
			fmt.Println(sortedParsedLines[i].toStringPlus())
		} else {
			fmt.Println(sortedParsedLines[i].toString())
		}
	}

}

func (l ParsedLineType) toStringPlus() string {
	return l.date + " " + l.hour + " " + l.timeZone + " " + l.pid + " " + internal.Green.String() + string(l.bddInfo) + " " + l.severityColor.String() + string(l.severity) + internal.Reset.String() + ":" + l.logMessage
}

func (l ParsedLineType) toString() string {
	return l.date + " " + l.hour + " " + l.timeZone + " " + l.pid + " " + l.severityColor.String() + string(l.severity) + internal.Reset.String() + ":" + l.logMessage
}

func sortLogsLines(lines []ParsedLineType, options Options) []ParsedLineType {

	// Start and End time parameters

	if options.StartTime != nil {
		//RECHERCHE DICHO !!!!
		//A OPTIMISER > parcours toutes les lignes et s'arrête quand atteint la date de départ
		i := 0
		for i < len(lines) {
			if lines[i].time.Compare(*options.StartTime) >= 0 {
				break
			} else {
				i++
			}
		}
		//A OPTIMISER > Si pas de lignes donc la boucle à tout parcouru (peut être su direct ...)
		if i == len(lines) {
			internal.PrintInfo("No line found after the start time " + options.StartTime.Format(time.UnixDate))
			return []ParsedLineType{}
		}
		lines = lines[i:]
	}

	if options.EndTime != nil {
		i := len(lines) - 1
		//Comme start date mais en partant de la fin, à optimiser pareil, on peut savoir si la première
		// ligne est après donc 0 lignes correcte, voir pour recherche dicho
		for i >= 0 {
			if lines[i].time.Compare(*options.EndTime) <= 0 {
				break
			} else {
				i--
			}
		}
		if i == -1 {
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

	//LogType parameter
	if options.LogType != ALL {
		typeLines := []ParsedLineType{}
		for i := range lines {
			if lines[i].logtype == options.LogType {
				typeLines = append(typeLines, lines[i])
			}
		}
		lines = typeLines
	}

	// NBLines paramater
	// possible de l'améliorer en retournant dès qu'on a le nombre
	// de ligne en fonction des options mais obligé d'être le dernier, à voir

	if len(lines) >= options.NBLines {
		lines = lines[(len(lines) - options.NBLines):]
	}

	return lines
}
