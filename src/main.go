package main

import (
	"fmt"
	Game "groupie/src/games"
)

func main() {
	var letters []string
	for i := 0; i < 26; i++ {
		letters = append(letters, Game.RandomLetter([]string{}))
		fmt.Println(letters[i])
	}

}
