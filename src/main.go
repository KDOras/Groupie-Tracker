package main

import (
<<<<<<< HEAD:Go/main.go
	"groupie/Go/databaseManager"
=======
	"fmt"
	"groupie/src/databaseManager"
>>>>>>> d4048eff6fdfa27bd9d4dac6ff220e0fc0a68c16:src/main.go
)

func main() {
	db := databaseManager.InitDatabase("SQL/database.db")

	databaseManager.ChangeRoomGameMode(db, 1, 2)

	// databaseManager.LeaveRoom(db, 2)

	// fmt.Println(err)
}
