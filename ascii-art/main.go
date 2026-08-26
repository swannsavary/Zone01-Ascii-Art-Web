package asciiart

import (
	"errors"
	"os"
	"strings"
)

// prend une chaine de caracere et renvoi une chaine de caractere en ascii-art
func AsciiArt(s string, typ string) (string, error) {
	fileContent, err := whoPolis(typ)
	if err != nil && err.Error() == "404" {
		return "", errors.New("404")
	} else if err != nil && err.Error() == "400" {
		return "", errors.New("400")
	}
	fileLines := strings.Split(string(fileContent), "\n")

	caractere, err := splitcara(s)
	if err != nil && err.Error() == "400" {
		return "", errors.New("400")
	}
	return changeToAsciiArt(fileLines, caractere), nil
}

// separe les caractere d'un string et le separe en un slice de slice de rune a chaque \n
func splitcara(s string) ([][]rune, error) {
	var result [][]rune
	s2 := []rune(s)
	s3 := []rune{}
	for i := 0; i < len(s2); i++ {
		if s2[i] >= 126 || s2[i] <= 32 {
			return [][]rune{}, errors.New("400") //<- return erreur 400 si un caractere n'est pas supporter avec le fichier ascii art
		}
		if i != len(s2)-1 && s2[i] == '\\' && s2[i+1] == 'n' {
			result = append(result, s3)
			s3 = []rune{}
		} else if i != 0 && s2[i-1] == '\\' && s2[i] == 'n' {
			continue
		} else {
			s3 = append(s3, s2[i])
		}
	}
	if len(s3) != 0 {
		result = append(result, s3)
	}
	return result, nil
}

// change un [][]rune en un autre en fonction de l'option de maniere qu'une fois afficher il l'affiche  en ascii-art
func changeToAsciiArt(file []string, cara [][]rune) string {
	var result string
	for j := 0; j < len(cara); j++ {
		if len(cara[j]) == 0 {
			result += "\n"
			continue
		}
		for i := 0; i <= 8; i++ {
			for e := 0; e < len(cara[j]); e++ {
				poslettre := (cara[j][e] - ' ') * 9
				result += file[i+int(poslettre)]
			}
			if i != 0 {
				result += "\n"
			}
		}
	}
	return result
}

// structure qui lie la valeur de l'option au fichier correspondant
var banners = map[string]string{
	"standard":   "standard.txt",
	"shadow":     "shadow.txt",
	"thinkertoy": "thinkertoy.txt",
}

// renvoie le bon fichier a exploiter en fonction de l'option actuelle
func whoPolis(typ string) ([]byte, error) {
	if chemin, exist := banners[typ]; exist {
		file, err := os.ReadFile(chemin)
		if err != nil {
			err = errors.New("404")
		}
		return file, err
	} else {
		err := errors.New("400")
		return nil, err
	}
}
