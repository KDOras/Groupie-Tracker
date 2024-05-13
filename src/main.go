package main

import (
	"groupie/src/databaseManager"
)

func main() {
	databaseManager.DelRoom(databaseManager.InitDatabase("SQL/database.db"), 2)
}
