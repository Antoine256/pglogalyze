package main

import (
	"fmt"
	"os"
	"pglogalyze/internal"
	"strings"
)

var file []string

func main() {

	//defines variables with parameters
	cmdParams := internal.ParseParameters(strings.Join(os.Args[1:], " "))
	if cmdParams == nil || len(*cmdParams) == 0 {
		return
	}

	params := *cmdParams
	first := params[0]

	fmt.Println(first.Param.Name)

	//get the file to read

}
