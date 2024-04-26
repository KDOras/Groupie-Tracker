package main

import (
	"groupie/src/databaseManager"
	Game "groupie/src/games"
)

func main() {
	db := databaseManager.InitDatabase("SQL/database.db")

	Game.AddPoint(db, 1)

	// databaseManager.LeaveRoom(db, 2)

	// fmt.Println(err)
}
