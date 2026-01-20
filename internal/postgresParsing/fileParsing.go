package postgresparsing

import (
	"regexp"
	"strings"
)

func parseFile(file []byte) []string {
	lines := []string{}
	stringFile := string(file)
	splitFile := strings.Split(stringFile, "\n")

	r, _ := regexp.Compile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d\d\d \w+ \[\d+\] \w+: .*`)

	i := 0
	buffer := ""
	for i < len(splitFile) {
		if i > 0 && r.MatchString(splitFile[i]) {
			//on valide la ligne précédente si la ligne actuelle commence par le log préfix
			lines = append(lines, buffer)
			buffer = ""
		}
		buffer += splitFile[i]
		i++
	}
	if r.MatchString(buffer) {
		lines = append(lines, buffer)
	}

	return lines
}
