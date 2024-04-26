package Game

import (
	"groupie/Go/databaseManager"
	"strings"
)

func VerifyUserResponse(user databaseManager.ConnectedUser, songName, userInput string) bool {
	songName = strings.ToLower(songName)
	userInput = strings.ToLower(userInput)

	if songName == userInput {
		return true
	}
	return false
}
