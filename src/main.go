package main

import (
	"fmt"
	"groupie/src/databaseManager"
)

func main() {
	fmt.Println(databaseManager.GetRoomFromUser(databaseManager.InitDatabase("SQL/database.db"), databaseManager.ConnectedUser{Id: 1}).Id)
}
