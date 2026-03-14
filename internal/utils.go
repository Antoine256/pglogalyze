package internal

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func StringToTime(strDate string, strHour string) time.Time {
	year, month, day := StringToTimeYMD(strDate)
	hour, minute, seconds := StringToTimeHMS(strHour)

	return time.Date(year, month, day, hour, minute, seconds, 0, time.UTC)
}

func StringToTimeYMD(str string) (int, time.Month, int) {

	splitStr := strings.Split(str, "-")

	year, err := strconv.Atoi(splitStr[0])
	intMonth, err := strconv.Atoi(splitStr[1])
	day, err := strconv.Atoi(splitStr[2])

	if err != nil {
		PrintError("Parsing line file StringToTimeYMD", err.Error())
	}

	month := time.Month(intMonth)

	return year, month, day

}

func StringToTimeHMS(str string) (int, int, int) {

	splitStr := strings.Split(str, ":")

	hour, err := strconv.Atoi(splitStr[0])
	minute, err := strconv.Atoi(splitStr[1])
	seconds, err := strconv.Atoi(strings.Split(splitStr[2], ".")[0])

	if err != nil {
		PrintError("Parsing line file StringToTimeHMS", err.Error())
	}

	return hour, minute, seconds
}
