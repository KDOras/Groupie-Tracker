package databaseManager

import (
	"database/sql"
	"log"

	"github.com/asaskevich/govalidator"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string
	Email    string
	Password string
}

type ConnectedUser struct {
	Id       int
	Username string
}

type Room struct {
	Id        int
	CreatedBy ConnectedUser
	MaxPlayer int
	Name      string
	Password  string
	GameMode  GameMode
}

type GameMode struct {
	Id   int
	Name string
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
					Username TEXT NOT NULL,
					email TEXT NOT NULL,
					password TEXT NOT NULL
				);
				
				CREATE TABLE IF NOT EXISTS ROOMS (
					id INTEGER PRIMARY KEY,
					created_by INTEGER NOT NULL,
					max_player INTEGER NOT NULL,
					name TEXT NOT NULL,
					password TEXT,
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

// Region Start - Rooms

func insertIntoRooms(db *sql.DB, room Room) (int64, error) {
	result, _ := db.Exec(`INSERT INTO ROOMS (created_by, max_player, name, password, id_game) VALUES (?, ?, ?, ?, ?)`, room.CreatedBy.Id, room.MaxPlayer, room.Name, room.Password, room.GameMode.Id)
	return result.LastInsertId()
}

func insertIntoRoom_Users(db *sql.DB, room Room, user ConnectedUser) (int64, error) {
	result, _ := db.Exec(`INSERT INTO ROOM_USERS (id_room, id_user) VALUES (?, ?)`, room.Id, user.Id)
	return result.LastInsertId()
}

func CreateRoom(db *sql.DB, room Room) {
	insertIntoRooms(db, room)
}

func GetRoom(db *sql.DB, id int) Room {
	data, err := db.Query(`SELECT * FROM ROOMS WHERE (id==?)`, id)
	if err != nil {
		log.Fatal(err)
	}
	var room Room
	var created_by int
	var game int
	for data.Next() {
		data.Scan(&room.Id, &created_by, room.MaxPlayer, room.Name, room.Password, &game)
	}
	room.CreatedBy = GetUserById(db, created_by)
	room.GameMode = getGame(db, game)
	return room
}

func JoinRoom(db *sql.DB, user ConnectedUser, room Room) string {
	if !isRoomFull(db, room) {
		insertIntoRoom_Users(db, room, user)
		return ""
	} else {
		return "This room is full."
	}
}

func isRoomFull(db *sql.DB, room Room) bool {
	data, err := db.Query(`SELECT id_user FROM ROOM_USERS WHERE (id_room==?)`, room)
	if err != nil {
		log.Fatal(err)
	}
	n := 0
	for data.Next() {
		n++
	}
	return n < room.MaxPlayer
}

func LeaveRoom(db *sql.DB, userId int) {
	db.Exec(`DELETE * FROM ROOM_USERS WHERE (id_user==?)`, userId)
}

func DelRoom(db *sql.DB, roomId int) {
	db.Exec(`DELETE * FROM ROOM_USERS WHERE (id_room==?)`, roomId)
}

func GetNumberOfPlayerFromRoom(db *sql.DB, roomId int) int {
	data, err := db.Query(`SELECT id_user FROM ROOM_USERS WHERE (id_room==?)`, roomId)
	if err != nil {
		log.Fatal(err)
	}
	n := 0
	for data.Next() {
		n++
	}
	return n
}

// Region End - Rooms

// Region Start - Game Modes

func getGame(db *sql.DB, id int) GameMode {
	data, err := db.Query(`SELECT * FROM GAMES WHERE (id==?)`, id)
	if err != nil {
		log.Fatal(err)
	}
	var game GameMode
	for data.Next() {
		data.Scan(&game.Id, &game.Name)
	}
	return game
}

// Region End - Game Modes

// Region Start - User

func LoggingIn(db *sql.DB, prompt string, password string) (ConnectedUser, string) {
	var user ConnectedUser
	var errBis ConnectedUser
	var pass string
	var data *sql.Rows
	var err error
	if govalidator.IsEmail(prompt) {
		data, err = db.Query(`SELECT id, Username, password FROM USER WHERE (email==?)`, prompt)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		data, err = db.Query(`SELECT id, Username, password FROM USER WHERE (Username==?)`, prompt)
		if err != nil {
			log.Fatal(err)
		}
	}
	for data.Next() {
		data.Scan(&user.Id, &user.Username, &pass)
	}
	if user.Username == "" {
		return errBis, "No user with this username found."
	}
	if pass == password {
		return user, ""
	}
	return errBis, "Incorrect password."
}

func GetUserById(db *sql.DB, id int) ConnectedUser {
	data, err := db.Query(`SELECT id, Username FROM USER WHERE (id==?)`, id)
	if err != nil {
		log.Fatal(err)
	}
	var user ConnectedUser
	for data.Next() {
		data.Scan(&user.Id, &user.Username)
	}
	return user
}

func GetUserByUsername(db *sql.DB, Username string) ConnectedUser {
	data, err := db.Query(`SELECT id, Username FROM USER WHERE (Username==?)`, Username)
	if err != nil {
		log.Fatal(err)
	}
	var user ConnectedUser
	for data.Next() {
		data.Scan(&user.Id, &user.Username)
	}
	return user
}

func GetUserByMail(db *sql.DB, email string) ConnectedUser {
	data, err := db.Query(`SELECT id, Username FROM USER WHERE (email==?)`, email)
	if err != nil {
		log.Fatal(err)
	}
	var user ConnectedUser
	for data.Next() {
		data.Scan(&user.Id, &user.Username)
	}
	return user
}

func insertIntoUser(db *sql.DB, user User) (int64, error) {
	result, _ := db.Exec(`INSERT INTO USER (Username, email, password) VALUES (?, ?, ?)`, user.Username, user.Email, user.Password)
	return result.LastInsertId()
}

func DeleteUserWithUsername(db *sql.DB, Username string) {
	db.Exec(`DELETE FROM USER WHERE (Username==?)`, Username)
}

func DeleteUserWithId(db *sql.DB, id int64) {
	db.Exec(`DELETE FROM USER WHERE (id==?)`, id)
}

func DeleteUserWithEmail(db *sql.DB, email string) {
	db.Exec(`DELETE FROM USER WHERE (email==?)`, email)
}

func CreateNewUser(db *sql.DB, user User) string {
	err := checkUsername(db, user.Username)
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

func checkUsername(db *sql.DB, Username string) string {
	rows, err := db.Query(`SELECT Username FROM USER WHERE (Username==?)`, Username)
	if err != nil {
		log.Fatal(err)
	}
	var user User
	for rows.Next() {
		err := rows.Scan(&user.Username)
		if err != nil {
			log.Fatal(err)
		}
	}
	if user.Username != "" {
		return "This Username is already used."
	}
	if len(Username) < 3 {
		return "Invalid Length, must be greater than 3."
	}
	for i := range Username {
		if string(Username[i]) == "@" {
			return "Invalid username, '@' is not a supported character"
		}
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
