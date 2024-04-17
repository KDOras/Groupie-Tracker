package main

import (
	"fmt"
	"groupie/Go/databaseManager"
)

func main() {
	db := databaseManager.InitDatabase("SQL/database.db")
	var user databaseManager.User
	user.Pseudo = "HoDoH"
	user.Email = "axelmichon.pro@gmail.com"
	user.Password = "Amemlm89260@"

	err := databaseManager.CreateNewUser(db, user)

	fmt.Println(err)
}
