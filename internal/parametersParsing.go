package internal

import (
	"os"
	"strconv"
	"strings"
)

type Parameter struct {
	Name        string
	Description string
	Strerror    string
	NbArgs      int
}

type CmdParam struct {
	Param Parameter
	Args  []string
}

func getParamsFromFile() []Parameter {

	parameters := []Parameter{}
	osfile, err := os.Open("./utils/parameters.csv")

	if err != nil {
		PrintError("internal\\parametersParsing.go", err.Error())
	}

	file, err := os.ReadFile(osfile.Name())

	if err != nil {
		PrintError("internal\\parametersParsing.go", err.Error())
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		splitLine := strings.Split(line, ";")
		i, err := strconv.Atoi(splitLine[3])
		if err != nil {
			PrintError("internal\\parametersParsing.go", err.Error())
		}
		parameters = append(parameters, Parameter{
			Name:        splitLine[0],
			Description: splitLine[1],
			Strerror:    splitLine[2],
			NbArgs:      i,
		})
	}

	return parameters
}

func getParam(name string, params []Parameter) *Parameter {
	found := false
	param := Parameter{}
	for _, p := range params {
		if p.Name == name {
			found = true
			param = p
		}
	}
	if !found {
		return nil
	}
	return &param
}

func ParseParameters(stringParams string) *[]CmdParam {
	params := getParamsFromFile()
	splitStringParams := strings.Split(stringParams, " -")
	cmdParams := []CmdParam{}

	for _, p := range splitStringParams {
		p = strings.TrimSpace(p)

		if len(p) == 0 {
			continue
		}
		splitp := strings.Split(p, " ")
		param := getParam(strings.Replace(splitp[0], "-", "", 1), params)
		if param == nil {
			paramsError(splitp[0], "not valid parameter")
			return nil
		}
		if len(splitp[1:]) != param.NbArgs {
			paramsError(splitp[0], "number of arguments is not valid, expected "+strconv.Itoa(param.NbArgs)+" args")
			return nil
		}
		cmdParams = append(cmdParams, CmdParam{
			Param: *param,
			Args:  splitp[1:],
		})
	}

	return &cmdParams
}

func GetParams(name string, params []CmdParam, numArgs int) string {
	res := ""
	for i := range params {
		if params[i].Param.Name == name {
			res = params[i].Args[numArgs]
		}
	}
	return res
}

func HasParams(name string, params []CmdParam) bool {
	res := false
	for i := range params {
		if name == params[i].Param.Name {
			res = true
		}
	}
	return res
}
