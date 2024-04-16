package main

import (
	"fmt"
	"groupie/Go/databaseManager"
)

func main() {
	db := databaseManager.InitDatabase("SQL/database.db")
	var user databaseManager.User
	user.Pseudo = "Mah"
	user.Email = "Test@MegaTest.com"
	user.Password = "Bruh"

	fmt.Println(db)
}
