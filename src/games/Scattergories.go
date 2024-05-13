package Game

import (
	"math/rand"
)

func RandomLetter(anciantLetters *[]string) string {
	letter := rand.Intn(26)
	formatLetter := string(byte(65 + letter))
	letterIsGood := true
	for _, i := range *anciantLetters {
		if formatLetter == string(i) {
			letterIsGood = false
		}
	}
	if letterIsGood {
		return formatLetter
	} else {
		return RandomLetter(anciantLetters)
	}
}
