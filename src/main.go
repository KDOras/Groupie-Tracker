package main

import (
	"fmt"
	"groupie/src/databaseManager"
)

func main() {
	db := databaseManager.InitDatabase("SQL/database.db")

	_, err := databaseManager.LoggingIn(db, "axelmichon.pro@gmail.com", "Amemlm89260@")
	fmt.Println(err)
}
