package main

import (
<<<<<<< HEAD
	"groupie/src/databaseManager"
=======
	"fmt"
	databaseManager "groupie/src/databasemanager"
>>>>>>> db02f1c8a46d943d706529c833d3cce82d8461ba
)

func main() {
	fmt.Println(databaseManager.GetRoomFromUser(databaseManager.InitDatabase("SQL/database.db"), databaseManager.ConnectedUser{Id: 1}).Id)
}
