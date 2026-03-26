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

	lines := getLines(osfile, options)

	fmt.Println("Lines found : " + internal.Yellow.String() + strconv.Itoa(len(lines)) + "/" + strconv.Itoa(options.NBLines) + internal.Reset.String())

	for i := 0; i < len(lines); i++ {
		//Permet d'aligner les lignes après le PID
		if len(lines[i].pid) < 7 {
			for j := 0; j < 7; j++ {
				if len(lines[i].pid) == 7 {
					break
				}
				lines[i].pid = lines[i].pid + " "
			}
		}
		if lines[i].bddInfo != "" {
			fmt.Println(lines[i].toStringPlus())
		} else {
			fmt.Println(lines[i].toString())
		}
	}

}

func (l ParsedLineType) toStringPlus() string {
	return l.date + " " + l.hour + " " + l.timeZone + " " + l.pid + " " + internal.Green.String() + string(l.bddInfo) + " " + l.severityColor.String() + string(l.severity) + internal.Reset.String() + ":" + l.logMessage
}

func (l ParsedLineType) toString() string {
	return l.date + " " + l.hour + " " + l.timeZone + " " + l.pid + " " + l.severityColor.String() + string(l.severity) + internal.Reset.String() + ":" + l.logMessage
}
