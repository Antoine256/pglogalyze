package main

import (
	"flag"
	"fmt"
	"pglogalyze/internal"
	postgresparsing "pglogalyze/internal/postgresParsing"
	"strconv"
	"strings"
)

var file = flag.String("f", "", "PostgreSQL log file")
var severityLevel = flag.String("l", "", "Severity level")
var start = flag.String("st", "", "Start time (YYYY-MM-DDTHH:MM:SS)")
var end = flag.String("et", "", "End time (YYYY-MM-DDTHH:MM:SS)")
var nbLines = flag.String("n", "", "Number of lines")

func main() {

	flag.Parse()

	options := postgresparsing.Options{LogFilePath: "", Level: postgresparsing.NONE}

	//---------------------- PARSING USER PARAMS ---------------------

	// PATH

	if *file != "" {
		path := *file
		if !internal.PathExists(path) {
			fmt.Println("Error: -f (log file) is not reachable")
			return
		} else {
			fmt.Println("Info: the logfile is : " + internal.Green + path + internal.Reset)
			options.LogFilePath = path
		}
	} else {
		internal.PrintInfo("Try to get log path by database informations")
		fmt.Println("Error: -f (log file) is required")
		//internal.GetPathByDatabaseConn(params)
		return
	}

	// SEVERITY

	if *severityLevel != "" {
		severity := *severityLevel
		if postgresparsing.IsAValidSeverity(severity) {
			options.Level = postgresparsing.Severity(severity)
		} else {
			fmt.Println("Error: -l (severity level) doesn t exist")
		}
	}

	// TIME

	if *start != "" {
		parseStartTime := strings.Split(*start, "T")
		strDate := parseStartTime[0]
		strHour := parseStartTime[1]
		time := internal.StringToTime(strDate, strHour)
		options.StartTime = &time
		fmt.Println("Info: Start date defined as : " + internal.Green + strDate + " " + strHour + internal.Reset)
	}

	if *end != "" {
		parseEndTime := strings.Split(*start, "T")
		strDate := parseEndTime[0]
		strHour := parseEndTime[1]
		time := internal.StringToTime(strDate, strHour)
		options.EndTime = &time
		fmt.Println("Info: End date defined as : " + internal.Green + strDate + " " + strHour + internal.Reset)
	}

	// Number of lines

	if *nbLines != "" {
		nb, err := strconv.Atoi(*nbLines)
		if err != nil {
			fmt.Println(err)
			return
		}
		options.NBLines = nb
	} else {
		options.NBLines = 30
	}
	fmt.Println("Info: Number of lines (default 30) : " + internal.Green + strconv.Itoa((options.NBLines)) + internal.Reset)

	//----------------------- READING LOG FILE -----------------------
	//fmt.Println(options)
	postgresparsing.ReadLogFile(options)

}
