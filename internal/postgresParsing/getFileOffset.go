package postgresparsing

import (
	"bufio"
	"errors"
	"os"
	"pglogalyze/internal"
	"regexp"
	"strings"
	"time"
)

func getTimeOffset(file *os.File, options Options, size int64) int64 {

	// On recherche avec dichotimie la une ligne dans les temps (au moins une lignes correspond)
	// on fait les saut de dicho, le but étant d'aller au plus près du End time donc si on est au milieur des deux on monte, (voir comment faire pour un curseur si c'est possible)

	offset := findOffset(file, *options.EndTime, size)

	return offset

}

func verifFirstAndLastLines(file *os.File, options Options) (bool, error) {
	//On teste la première ligne du fichier (si elle est postérieur à endTime alors pas de lignes trouvées)
	if options.EndTime != nil {
		scanner := bufio.NewScanner(file)

		if scanner.Scan() {
			line := scanner.Text()
			if strings.Trim(line, " ") != "" {
				parsedLine := parseLine(line)
				// Si la première ligne est supérieur au filtre end time alors aucune ligne ne correspond
				if parsedLine.time.Compare(*options.EndTime) > 0 {
					return false, errors.New("First line before End time !")
				}
			}
		}
	}

	// On teste la dernière ligne du fichier (si elle est antérieur à startTime alors pas de lignes trouvées)
	if options.StartTime != nil {
		line := getLastLineOfFile(*file)
		parsedLine := parseLine(line)
		// Si la dernière ligne est avant le starttime alors aucune ligne du fichier ne corerspond au filtre
		if parsedLine.time.Compare(*options.StartTime) < 0 {
			return false, errors.New("Last line after Start time !")
		}
	}

	return true, nil
}

func isValidLineOmitTime(line ParsedLineType, options Options) bool {

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

func getLastLineOfFile(file os.File) string {

	r, _ := regexp.Compile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d\d\d \w+ \[\d+\] \w+: .*`)
	r2, _ := regexp.Compile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d\d\d \w+ \[\d+\] \w+@\w+ \w+: .*`)
	buffer := ""

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
			if r.MatchString(line) || r2.MatchString(line) {
				return line
			}
		}
	}

	return ""
}

func findOffset(file *os.File, target time.Time, size int64) int64 {
	r := regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d\d\d \w+ \[\d+\]`)
	low := int64(0)
	high := size
	found := false
	const blockSize = 4096
	const maxBufferSize = 1 * 1024 * 1024 // 1MB
	result := int64(0)
	//resultLine := ""
	for !found {
		mid := (low + high) / 2
		// Chercher deux timestamps consécutifs
		timestamps := []time.Time{}

		// on avance tant qu'on a pas de time de log, si on arrive à la fin alors on repars de mid vers l'arrière. ( toujours prendre deux lignes qui se suivent et regarder si nitre date est entre (dans ce cas on sort), après on met low à mid et avant on met hight à mid)
		tempMid := mid

		reachedEOF := false
		remainder := ""
		for len(timestamps) != 2 {
			buf := make([]byte, blockSize)
			file.ReadAt(buf, tempMid)
			stringBuf := ""
			lines := []string{}
			stringBuf = remainder + string(buf)
			if len([]byte(stringBuf)) > maxBufferSize {
				panic("buffer too large — probablement un problème de parsing")
			}
			lines = strings.Split(stringBuf, "\n")
			remainder = lines[len(lines)-1]

			// Chercher la première ligne qui matche un préfixe de log
			tempOffset := tempMid + int64(len([]byte(remainder)))
			for _, line := range lines[1:] {
				tempOffset += int64(len([]byte(line)))
				if r.MatchString(line) {
					timestamps = append(timestamps, getTimeOfLine(line))
					if len(timestamps) == 2 {
						result = tempOffset
						//resultLine = line
						break
					}
				}
			}

			tempMid += blockSize
			if tempMid > size {
				if reachedEOF {
					break
				}
				tempMid = size
				reachedEOF = true
			}
		}

		//Je nai pas trouvé deux timestamp au dessus, je regarde en dessous
		reached0 := false
		remainder = ""
		tempMid = mid - blockSize
		for len(timestamps) != 2 {
			buf := make([]byte, blockSize)
			file.ReadAt(buf, tempMid)

			stringBuf := ""
			lines := []string{}
			stringBuf = string(buf) + remainder
			lines = strings.Split(stringBuf, "\n")
			if len([]byte(stringBuf)) > maxBufferSize {
				panic("buffer too large — probablement un problème de parsing")
			}
			remainder = lines[0]

			// Chercher la première ligne qui matche un préfixe de log en partant de la fin
			tempOffset := tempMid + blockSize
			for i := len(lines) - 1; i >= 1; i-- {
				line := lines[i]
				tempOffset -= int64(len([]byte(line)))
				if r.MatchString(line) {
					timestamps = append(timestamps, getTimeOfLine(line))
					if len(timestamps) == 2 {
						result = tempOffset
						//resultLine = line
						break
					}
				}
			}

			tempMid -= blockSize
			if tempMid < 0 {
				if reached0 {
					break
				}
				tempMid = 0
				reached0 = true
			}
		}

		// we have the two timestamp, verify the conditions and moove borders
		if len(timestamps) != 2 {
			return 0
		} else if timestamps[0].After(timestamps[1]) {
			timestamps[0], timestamps[1] = timestamps[1], timestamps[0]
		}

		if timestamps[0].After(target) {
			high = mid
		} else if timestamps[1].After(target) {
			found = true
		} else {
			low = mid
		}
	}

	//fmt.Fprintln(os.Stdin, internal.Blue, "Result Line :", resultLine, internal.Reset)

	return result
}

func getTimeOfLine(line string) time.Time {
	line = strings.TrimSpace(line)
	splitLine := strings.Split(line, " ")
	lineDate := splitLine[0]
	lineHour := splitLine[1]
	return internal.StringToTime(lineDate, lineHour)
}
