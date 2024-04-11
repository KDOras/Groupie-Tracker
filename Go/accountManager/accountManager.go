package accountManager

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string
	Password string
	Email    string
	Id       int
}

func AddNewAccount(newUser User) {
	// Add new user to database
}

func LoggingIn(username string, password string) {
	db, err := sql.Open("mysql", "/database.sql")
	checkErr(err)
	defer db.Close()

	insert := `INSERT INTO USER(id, pseudo, email, password) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insert)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(0, "test", "test", "test")
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
