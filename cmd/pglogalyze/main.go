package main

import (
	"fmt"
	"os"
	"pglogalyze/internal"
	postgresparsing "pglogalyze/internal/postgresParsing"
	"strings"
)

var file []string

func main() {

	options := postgresparsing.Options{LogFilePath: "", Level: postgresparsing.NONE}

	//---------------------- PARSING USER PARAMS ---------------------

	cmdParams := internal.ParseParameters(strings.Join(os.Args[1:], " "))
	if cmdParams == nil {
		return
	}

	params := *cmdParams

	// PATH

	if internal.HasParams("f", params) {
		path := internal.GetParams("f", params, 0)
		if !internal.PathExists(path) {
			internal.PrintError("./main", "specified file doesn't exist : "+path)
		} else {
			options.LogFilePath = path
		}
	} else {
		internal.PrintInfo("Try to get log path by database informations")
		//internal.GetPathByDatabaseConn(params)
	}

	// SEVERITY

	if internal.HasParams("l", params) {
		severity := internal.GetParams("l", params, 0)
		if postgresparsing.IsAValidSeverity(severity) {
			options.Level = postgresparsing.Severity(severity)
		} else {
			internal.PrintError("./main", "Level of severity is not correct : "+severity)
		}
	}

	// TIME

	if internal.HasParams("st", params) {
		strDate := internal.GetParams("st", params, 0)
		strHour := internal.GetParams("st", params, 1)

		time := internal.StringToTime(strDate, strHour)

		options.StartTime = &time
	}

	if internal.HasParams("et", params) {
		strDate := internal.GetParams("et", params, 0)
		strHour := internal.GetParams("et", params, 1)

		time := internal.StringToTime(strDate, strHour)

		options.EndTime = &time
	}

	//----------------------- READING LOG FILE -----------------------
	fmt.Println(options)
	postgresparsing.ReadLogFile(options)

}
