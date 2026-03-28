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

	lines := make([]ParsedLineType, 0, options.NBLines)
	stat, _ := osfile.Stat()
	size := stat.Size()

	// If time options, verify first and last line...
	if options.StartTime != nil || options.EndTime != nil {
		verifLines, err := verifFirstAndLastLines(osfile, options)
		if !verifLines {
			if err != nil {
				fmt.Fprintln(os.Stderr, internal.Red, "ERROR ", err.Error())
			}
		}
	}

	// if Endtime, get offset closest line to endtime
	if options.EndTime != nil {
		if parseLine(getLastLineOfFile(*osfile)).time.Compare(*options.EndTime) < 0 {
			fmt.Fprintln(os.Stdout, internal.Blue, "INFO : La dernière ligne est antérieur au paramètre EndTime")
		} else {
			offset := getTimeOffset(osfile, options, size)
			size = offset
		}
	}

	getLines(osfile, options, &lines, size)

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
