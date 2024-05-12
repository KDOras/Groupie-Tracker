package main

import (
	databaseManager "groupie/src/databasemanager"
)

func main() {
	databaseManager.DelRoom(databaseManager.InitDatabase("SQL/database.db"), 2)
}
