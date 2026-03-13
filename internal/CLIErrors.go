package internal

import (
	"fmt"
	"os"
)

type CMDCOLOR int

const (
	Red CMDCOLOR = iota
	Yellow
	Green
	Blue
	Reset
)

var colorstring = map[CMDCOLOR]string{
	Red:    "\033[1;91m",
	Yellow: "\033[1;93m",
	Green:  "\033[1;92m",
	Blue:   "\033[1;94m",
	Reset:  "\033[0m",
}

func (cc CMDCOLOR) String() string {
	if s, ok := colorstring[cc]; ok {
		return s
	}
	return ""
}

func paramsError(parameterName string, err string) {
	fmt.Fprintln(os.Stderr, "wrong parameter ("+parameterName+") :", Red, err, Reset)
}

func PrintError(file string, err string) {
	fmt.Fprintln(os.Stderr, "error in file "+file+" : ", Red, err, Reset)
}

func PrintInfo(info string) {
	fmt.Fprintln(os.Stderr, Blue.String()+"info : ", info, Reset)
}
