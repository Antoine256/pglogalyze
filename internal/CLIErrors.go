package internal

import (
	"fmt"
	"os"
)

const (
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Reset  = "\033[0m"
)

func paramsError(parameterName string, err string) {
	fmt.Fprintln(os.Stderr, Red+"wrong parameter ("+parameterName+") :"+Reset, err)
}

func internError(file string, err string) {
	fmt.Fprintln(os.Stderr, Red+"intern error file "+file+" : "+Reset, err)
}
