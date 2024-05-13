package Game

import (
	"fmt"
	"groupie/src/databaseManager"
	"math/rand"
	"os"
	"strings"
)

type Song struct {
	Name   string
	Lyrics string
}

func VerifyUserResponse(songName, userInput string) bool {
	songName = strings.ToLower(songName)
	userInput = strings.ToLower(userInput)

	if songName == userInput {
		return true
	}
	return false
}

func GetHighestId() int {
	db := databaseManager.InitDatabase("SQL/database.db")
	data, _ := db.Query("SELECT id from Songs ORDER BY id DESC LIMIT 1")
	id := 0
	for data.Next() {
		data.Scan(&id)
	}
	return id
}

func GetRandomLyrics() Song {
	db := databaseManager.InitDatabase("SQL/database.db")
	n := rand.Intn(GetHighestId())
	data, err := db.Query("SELECT name, lyrics FROM Songs WHERE id=?", n)
	if err != nil {
		fmt.Println(err)
		return Song{}
	}
	var lyrics string
	var name string
	for data.Next() {
		data.Scan(&name, &lyrics)
	}
	song := Song{Name: name, Lyrics: lyrics}
	return song
}

func GetRandomSong() string {
	dir, _ := os.ReadDir("./../mp3")
	n := rand.Intn(len(dir) - 1)
	str := fmt.Sprintf("%v", dir[n])
	return str
}
