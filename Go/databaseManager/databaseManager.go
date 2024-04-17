package databaseManager

import (
	"database/sql"
	"log"

	"github.com/asaskevich/govalidator"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Pseudo   string
	Email    string
	Password string
}

// Region Start - Database

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

// Region End - Database

// Region Start - User

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

func CreateNewUser(db *sql.DB, user User) string {
	err := checkPseudo(db, user.Pseudo)
	if len(err) != 0 {
		return err
	}
	err = checkMail(db, user.Email)
	if len(err) != 0 {
		return err
	}
	err = checkPass(user.Password)
	if len(err) != 0 {
		return err
	}
	insertIntoUser(db, user)
	return ""
}

func checkPseudo(db *sql.DB, pseudo string) string {
	rows, err := db.Query(`SELECT pseudo FROM USER WHERE (pseudo==?)`, pseudo)
	if err != nil {
		log.Fatal(err)
	}
	var user User
	for rows.Next() {
		err := rows.Scan(&user.Pseudo)
		if err != nil {
			log.Fatal(err)
		}
	}
	if user.Pseudo != "" {
		return "This pseudo is already used."
	}
	if len(pseudo) < 3 {
		return "Invalid Length, must be greater than 3."
	}
	return ""
}

func checkMail(db *sql.DB, email string) string {
	if !govalidator.IsEmail(email) {
		return "Not a valid adress."
	} else {
		rows, err := db.Query(`SELECT email FROM USER WHERE (email==?)`, email)
		if err != nil {
			log.Fatal(err)
		}
		var user User
		for rows.Next() {
			err = rows.Scan(&user.Email)
			if err != nil {
				log.Fatal(err)
			}
		}
		if user.Email != "" {
			return "This mail address is already used."
		}
	}
	return ""
}

func checkPass(password string) string {
	if len(password) < 12 {
		return "This password is too short, must be atleast 12 characters."
	}
	gotUpper := false
	gotLower := false
	gotNumber := false
	gotSpecial := false
	for _, e := range password {
		if e >= 65 && e <= 90 && !gotUpper {
			gotUpper = true
		}
		if e >= 97 && e <= 122 && !gotLower {
			gotLower = true
		}
		if e >= 48 && e <= 57 && !gotNumber {
			gotNumber = true
		}
		if e >= 33 && e <= 47 && !gotSpecial || e >= 58 && e <= 64 && !gotSpecial || e >= 91 && e <= 96 && !gotSpecial || e >= 123 && e <= 126 && !gotSpecial {
			gotSpecial = true
		}
	}
	if !gotUpper {
		return "This password doesn't contain uppercase."
	}
	if !gotLower {
		return "This password doesn't contain lowercase."
	}
	if !gotNumber {
		return "This password doesn't contain number."
	}
	if !gotSpecial {
		return "This password doesn't contain special character."
	}
	return ""
}

// Region End - User
