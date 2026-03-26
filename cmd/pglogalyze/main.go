package main

import (
	"flag"
	"fmt"
	"os"
	"pglogalyze/internal"
	postgresparsing "pglogalyze/internal/postgresParsing"
	"strconv"
	"strings"
	"time"
)

var file = flag.String("f", "", "PostgreSQL log file")
var severityLevel = flag.String("l", "", "Severity level")
var logType = flag.String("t", "", "Type of the log : QUERY | CONNECTION | DURATION | CHECKPOINT | STARTUP | SHUTDOWN")
var start = flag.String("st", "", "Start time (YYYY-MM-DDTHH:MM:SS)")
var end = flag.String("et", "", "End time (YYYY-MM-DDTHH:MM:SS)")
var nbLines = flag.Int("n", 20, "Number of lines")

func main() {

	timeStart := time.Now()

	flag.Parse()

	options := postgresparsing.Options{LogFilePath: "", Level: postgresparsing.NONE}

	//---------------------- PARSING USER PARAMS ---------------------

	// PATH

	if *file != "" {
		path := *file
		if !internal.PathExists(path) {
			fmt.Fprintln(os.Stdout, internal.Red, "Error", internal.Reset, ": -f (log file) is not reachable")
			return
		} else {
			fmt.Fprintln(os.Stdout, "the logfile is :", internal.Yellow, path, internal.Reset)
			options.LogFilePath = path
		}
	} else {
		fmt.Fprintln(os.Stderr, internal.Red, "Error", internal.Reset, ": -f (log file) is required")
		return
	}

	// SEVERITY

	if *severityLevel != "" {
		severity := *severityLevel
		if postgresparsing.IsAValidSeverity(severity) {
			options.Level = postgresparsing.Severity(severity)
		} else {
			fmt.Fprintln(os.Stderr, internal.Red, "Error", internal.Reset, ": -l (severity level) doesn t exist")
		}
	}

	// LOGTYPE

	if *logType != "" {
		LType := strings.ToUpper(*logType)
		if postgresparsing.IsAValidType(LType) {
			options.LogType = postgresparsing.LType(LType)
		} else {
			fmt.Fprintln(os.Stderr, internal.Red, "Error", internal.Reset, ": -l (severity level) doesn t exist")
		}
	}

	// TIME

	if *start != "" {
		parseStartTime := strings.Split(*start, "T")
		if len(parseStartTime) != 2 {
			fmt.Fprintln(os.Stderr, internal.Red, "Error", internal.Reset, ": -st format is not good")
			return
		}
		strDate := parseStartTime[0]
		strHour := parseStartTime[1]
		time := internal.StringToTime(strDate, strHour)
		options.StartTime = &time
		fmt.Fprintln(os.Stdout, "Start time defined as :", internal.Yellow, strDate, strHour, internal.Reset)
	}

	if *end != "" {
		parseEndTime := strings.Split(*end, "T")
		if len(parseEndTime) != 2 {
			fmt.Fprintln(os.Stderr, internal.Red, "Error", internal.Reset, ": -et format is not good")
			return
		}
		strDate := parseEndTime[0]
		strHour := parseEndTime[1]
		time := internal.StringToTime(strDate, strHour)
		options.EndTime = &time
		fmt.Fprintln(os.Stdout, "End time defined as :", internal.Yellow, strDate, strHour, internal.Reset)
	}

	// Number of lines

	if *nbLines != 0 {
		nb := *nbLines
		options.NBLines = nb
	} else {
		fmt.Fprintln(os.Stderr, internal.Red, "Error", internal.Reset, ": -n = 0")
		return
	}

	fmt.Fprintln(os.Stdout, "Number of lines : ", internal.Yellow, strconv.Itoa((options.NBLines)), internal.Reset)

	//----------------------- READING LOG FILE -----------------------
	fmt.Println("---------------------- LOGS ----------------------")
	postgresparsing.ReadLogFile(options)

	elapsed := time.Since(timeStart)
	fmt.Printf("Resolved in : %s\n", elapsed)

}
