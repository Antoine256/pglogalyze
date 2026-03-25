package postgresparsing

import (
	"os"
	"regexp"
	"strings"
)

// ////
//
// Fonction pour lire et récupérer les lignes du fichier de log en fonction des paramètre passé par l'utilisateur.
// La fonction lit le fichier en partant de la fin, block de bits par block de bits
// Chaque block est annalysé au fur et à mesure, les lignes sont reconnues et placé dans un tableau,
// le tableau est parcouru, si la lignes correspond à un début de log elle est parsée,
// et si elle correpond aux paramètres en entré elle est ajoutée, la fonction se stop lorsque le nombre de ligne voulu est atteint.

func getLines(file *os.File, options Options) []ParsedLineType {
	r, _ := regexp.Compile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d\d\d \w+ \[\d+\] \w+: .*`)
	r2, _ := regexp.Compile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d\d\d \w+ \[\d+\] \w+@\w+ \w+: .*`)
	buffer := ""
	lines := []ParsedLineType{}

	const blockSize = 4096
	const maxBufferSize = 1 * 1024 * 1024 // 1MB
	var remainder []byte

	stat, _ := file.Stat()
	size := stat.Size()

	for offset := size; offset > 0; {
		// Initialisation de readSize (taille du bloc lu)
		readSize := blockSize
		if offset < blockSize {
			readSize = int(offset)
		}

		// On soustrait le nombre de bits lu au "curseur" et on créer un buffeur de cette taille
		offset -= int64(readSize)
		buf := make([]byte, readSize)

		// func (f *File) ReadAt(b []byte, off int64) (n int, err error)
		// ReadAt reads len(b) bytes from the File starting at byte offset off. It returns the number of bytes read and the error, if any.
		// ReadAt always returns a non-nil error when n < len(b). At end of file, that error is io.EOF.
		// En gros on lit les bit à partir de notre curseur (offset) pour remplir notre buffeur.
		file.ReadAt(buf, offset)

		// On ajoute les bits restant
		buf = append(buf, remainder...)
		// Et on transforme en string
		stringBuf := string(buf)

		// Limite de la taille du buffer pour les cas ou le parsing ne fonctionne pas.
		if len(buf) > maxBufferSize {
			panic("buffer too large — probablement un problème de parsing")
		}

		// On découpe notre buf (avec le remainder) en fonction des retour à la lignes (voir si possible de mettre direct regex)
		splitBufByLines := strings.Split(stringBuf, "\n")
		// On met ce qui ne formait pas une ligne dans le remainder.
		remainder = []byte(splitBufByLines[0])

		// boucle sur les lignes dans l'ordre inverse
		for i := len(splitBufByLines) - 1; i >= 1; i-- {
			line := string(splitBufByLines[i])
			buffer = line + buffer
			// fmt.Println(buffer)

			if r2.MatchString(line) {
				// Si c'est une ligne appli, on ajoute avec l'indicateur
				buffer = "[\\]" + buffer
			}
			if r.MatchString(line) || r2.MatchString(line) {
				// Si la ligne corerespond à un début de log, on vérifie qu'elle correspond aux paramètres
				parsedLine := parseLine(buffer)
				if isValidLine(parsedLine, options) {
					lines = append(lines, parsedLine)
				}
				buffer = ""
			}
		}

		if len(lines) >= options.NBLines {
			return lines[:options.NBLines]
		}

	}

	// GERER CE QUIL RESTE DANS LE REMAINDER... première ligne du fichier normalement
	if len(remainder) > 0 {
		buffer = string(remainder) + buffer
		// traiter buffer comme une dernière ligne
		if r2.MatchString(string(remainder)) {
			buffer = "[\\]" + buffer
		}
		if r.MatchString(string(remainder)) || r2.MatchString(string(remainder)) {
			// Si la ligne corerespond à un début de log, on vérifie qu'elle correspond aux paramètres
			parsedLine := parseLine(buffer)
			if isValidLine(parsedLine, options) {
				lines = append(lines, parsedLine)
			}
			buffer = ""
		}
	}

	// Vérfier les params et ajouter la lignes (si ici alors encore de la place sinon return plus hiat dès que possible)

	return lines

}

func isValidLine(line ParsedLineType, options Options) bool {
	// Start and End time parameters

	if options.StartTime != nil {
		if line.time.Compare(*options.StartTime) < 0 {
			return false
		}
	}

	if options.EndTime != nil {
		if line.time.Compare(*options.EndTime) > 0 {
			return false
		}
	}

	// Severity parameter

	if options.Level != NONE {
		if line.severity != options.Level {
			return false
		}
	}

	//LogType parameter
	if options.LogType != ALL {
		if line.logtype != options.LogType {
			return false
		}
	}

	return true
}
