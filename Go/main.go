package main

import (
	"groupie/Go/databaseManager"
)

func main() {
	db := databaseManager.InitDatabase("SQL/database.db")

	databaseManager.ChangeRoomGameMode(db, 1, 2)

	// databaseManager.LeaveRoom(db, 2)

	// fmt.Println(err)
}
