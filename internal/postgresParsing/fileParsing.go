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
	r2, _ := regexp.Compile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d\d\d \w+ \[\d+\] \w+@\w+ \w+: .*`)
	i := 0
	buffer := ""
	bufferAP := false
	for i < len(splitFile) {
		if i > 0 && r.MatchString(splitFile[i]) || r2.MatchString(splitFile[i]) {
			//on valide la ligne précédente si la ligne actuelle commence par le log préfix
			if bufferAP {
				//si c'était une ligne appli, on ajoute avec l'indicateur
				buffer = "[\\]" + buffer
			}
			lines = append(lines, buffer)
			buffer = ""
			bufferAP = false
		}
		buffer += splitFile[i]
		if r2.MatchString(splitFile[i]) {
			// c'est une ligne avec un user@database (il faut pouvoir les reconnaitre)
			bufferAP = true
		}
		i++
	}
	if r.MatchString(buffer) {
		lines = append(lines, buffer)
	}

	return lines
}
