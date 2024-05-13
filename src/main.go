package main

import (
	"fmt"
	databaseManager "groupie/src/databasemanager"
)

func main() {
	fmt.Println(databaseManager.GetRoomFromUser(databaseManager.InitDatabase("SQL/database.db"), databaseManager.ConnectedUser{Id: 1}).Id)
}
