package databaseManager

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Pseudo   string
	Email    string
	Password string
}

func InitDatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)

	if err != nil {
		log.Fatal(err)
	}

	sqltStmt := `
				CREATE TABLE IF NOT EXISTS USER (
					id INTEGER PRIMARY KEY,
					pseudo TEXT NOT NULL,
					email TEXT NOT NULL,
					password TEXT NOT NULL
				);
				
				CREATE TABLE IF NOT EXISTS ROOMS (
					id INTEGER PRIMARY KEY,
					created_by INTEGER NOT NULL,
					max_player INTEGER NOT NULL,
					name TEXT NOT NULL,
					id_game INTEGER,
					FOREIGN KEY (created_by) REFERENCES USER(id),
					FOREIGN KEY (id_game) REFERENCES GAMES(id)
				);
				
				CREATE TABLE IF NOT EXISTS ROOM_USERS (
					id_room INTEGER,
					id_user INTEGER,
					score INTEGER,
					FOREIGN KEY (id_room) REFERENCES ROOMS(id),
					FOREIGN KEY (id_user) REFERENCES USER(id),
					PRIMARY KEY (id_room, id_user)
				);
				
				CREATE TABLE IF NOT EXISTS GAMES (
					id INTEGER PRIMARY KEY,
					name TEXT NOT NULL
				);
				`

	_, err = db.Exec(sqltStmt)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func insertIntoUser(db *sql.DB, user User) (int64, error) {
	result, _ := db.Exec(`INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)`, user.Pseudo, user.Email, user.Password)
	return result.LastInsertId()
}

func DeleteUserWithPseudo(db *sql.DB, pseudo string) {
	db.Exec(`DELETE FROM USER WHERE (pseudo==?)`, pseudo)
}

func DeleteUserWithId(db *sql.DB, id int64) {
	db.Exec(`DELETE FROM USER WHERE (id==?)`, id)
}

func DeleteUserWithEmail(db *sql.DB, email string) {
	db.Exec(`DELETE FROM USER WHERE (email==?)`, email)
}

func CreateNewUser(db *sql.DB, pseudo, email, password string) string {
	var user User
	err := checkPseudo(pseudo)
	if len(err) != 0 {
		return err
	}
	err = checkMail(email)
	if len(err) != 0 {
		return err
	}
	err = checkPass(password)
	if len(err) != 0 {
		return err
	}
	user.Pseudo, user.Email, user.Password = pseudo, email, password
	return ""
}

func checkPseudo(pseudo string) string {
	var err string
	if len(pseudo) < 3 {
		err = "Invalid Length, must be greater than 3"
	}
	return err
}

func checkMail(pseudo string) string {
	var err string
	return err
}

func checkPass(pseudo string) string {
	var err string
	return err
}
