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
	//On va à la première ligne correspondante à un log :

	for i < len(splitFile) && (r.MatchString(splitFile[i]) == false && r2.MatchString(splitFile[i]) == false) {
		i++
	}

	buffer := splitFile[i]
	bufferAP := false

	if r2.MatchString(splitFile[i]) {
		// c'est une ligne avec un user@database (il faut pouvoir les reconnaitre)
		bufferAP = true
	}
	for i < len(splitFile) {
		if i > 0 && (r.MatchString(splitFile[i]) || r2.MatchString(splitFile[i])) {
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
