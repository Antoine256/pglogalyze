package internal

import (
	"fmt"
	"os"
)

const (
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

func paramsError(parameterName string, err string) {
	fmt.Fprintln(os.Stderr, "wrong parameter ("+parameterName+") :"+Red, err, Reset)
}

func PrintError(file string, err string) {
	fmt.Fprintln(os.Stderr, "error in file "+file+" : "+Red, err, Reset)
}

func PrintInfo(info string) {
	fmt.Fprintln(os.Stderr, Blue+"info : ", info, Reset)
}
