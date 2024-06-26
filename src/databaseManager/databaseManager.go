package databaseManager

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"strconv"

	emailverifier "github.com/AfterShip/email-verifier"
	_ "github.com/mattn/go-sqlite3"
)

var (
	verifier = emailverifier.NewVerifier()
)

type User struct {
	Username string
	Email    string
	Password string
}

type ConnectedUser struct {
	Id       interface{}
	Username interface{}
}

type PlayerScore struct {
	User  string
	Score int
}

type LeaderBoard struct {
	UserList []PlayerScore
}

type Room struct {
	Id             int
	CreatedBy      ConnectedUser
	MaxPlayer      int
	NumberOfPlayer int
	Name           string
	Password       string
	GameMode       GameMode
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
	res, _ := db.Exec(`INSERT INTO ROOMS (created_by, max_player, name, password, id_game) VALUES (?, ?, ?, ?, ?)`, room.CreatedBy.Id, room.MaxPlayer, room.Name, room.Password, room.GameMode.Id)
	return res.LastInsertId()
}

func insertIntoRoom_Users(db *sql.DB, room int64, user ConnectedUser) {
	db.Exec(`INSERT INTO ROOM_USERS (id_room, id_user, score) VALUES (?, ?, ?)`, room, user.Id, 0)
}

func insertIntoLeaderboard(db *sql.DB, room int64) {
	db.Exec(`INSERT INTO LEADERBOARDS (room_id, first, second, third, fourth, fifth) VALUES (?, ?, ?, ?, ?, ?)`, room, "", "", "", "", "")
}

func ChangeRoomGameMode(db *sql.DB, idRoom, idGame int) {
	db.Exec(`UPDATE ROOMS SET id_game=? WHERE id=?`, idGame, idRoom)
}

func CreateRoom(db *sql.DB, room Room) {
	id, _ := insertIntoRooms(db, room)
	defer insertIntoRoom_Users(db, id, room.CreatedBy)
	defer insertIntoLeaderboard(db, id)
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
		data.Scan(&room.Id, &created_by, &room.MaxPlayer, &room.Name, &room.Password, &game)
	}
	room.CreatedBy = GetUserById(db, created_by)
	room.GameMode = GetGame(db, game)
	return room
}

func GetRoomFromUser(db *sql.DB, user ConnectedUser) Room {
	data, err := db.Query(`SELECT id_room FROM ROOM_USERS WHERE (id_user==?)`, user.Id)
	if err != nil {
		log.Fatal(err)
	}
	id := 0
	for data.Next() {
		data.Scan(&id)
	}
	return GetRoom(db, id)
}

func JoinRoom(db *sql.DB, user ConnectedUser, room Room) string {
	if !isRoomFull(db, room) {
		if !IsAlreadyPlaying(db, user) {
			insertIntoRoom_Users(db, int64(room.Id), user)
			return ""
		}
		return "Someone with the same account is already playing."
	} else {
		return "This room is full."
	}
}

func isRoomFull(db *sql.DB, room Room) bool {
	data, err := db.Query(`SELECT id_user FROM ROOM_USERS WHERE (id_room==?)`, room.Id)
	if err != nil {
		log.Fatal(err)
	}
	n := 0
	for data.Next() {
		n++
	}
	return n >= room.MaxPlayer
}

func giveLead(db *sql.DB, roomId int, actualLead int) (int64, error) {
	data, _ := db.Query(`SELECT id_user FROM ROOM_USERS WHERE id_user!=? AND id_room=?`, actualLead, roomId)
	userList := []int{}
	for data.Next() {
		var id int
		data.Scan(&id)
		userList = append(userList, id)
	}
	if len(userList) > 1 {
		result, _ := db.Exec(`UPDATE ROOMS SET created_by=? WHERE id=?`, userList[rand.Intn(len(userList)-1)], roomId)
		return result.LastInsertId()
	} else {
		result, _ := db.Exec(`UPDATE ROOMS SET created_by=? WHERE id=?`, userList[0], roomId)
		return result.LastInsertId()
	}
}

func LeaveRoom(db *sql.DB, userId int) {
	data, _ := db.Query(`SELECT id, created_by FROM ROOMS WHERE id=(SELECT id_room FROM ROOM_USERS WHERE id_user=?)`, userId)
	defer db.Close()
	var created_by int
	var id int
	for data.Next() {
		data.Scan(&id, &created_by)
	}
	if created_by == userId && GetNumberOfPlayerFromRoom(db, id) > 1 {
		giveLead(db, id, created_by)
	} else if created_by == userId && GetNumberOfPlayerFromRoom(db, id) == 1 {
		DelRoom(db, id)
	}
	db.Exec(`DELETE FROM ROOM_USERS WHERE id_user=?`, userId)
}

func DelRoom(db *sql.DB, roomId int) {
	db.Exec(`DELETE FROM ROOM_USERS WHERE (id_room==?)`, roomId)
	db.Exec(`DELETE FROM ROOMS WHERE (id==?)`, roomId)
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

func GetGame(db *sql.DB, id int) GameMode {
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

func IncreaseUserScore(user ConnectedUser) {
	db := InitDatabase("SQL/database.db")
	data, _ := db.Query(`SELECT score FROM ROOM_USERS WHERE (id_user==?)`, user.Id)
	score := 0
	for data.Next() {
		data.Scan(&score)
	}
	score++
	db.Exec("UPDATE ROOM_USERS SET score=? WHERE id_user=?", score, user.Id)
}

func ResetUserScore(user ConnectedUser) {
	db := InitDatabase("SQL/database.db")
	db.Exec("UPDATE ROOM_USERS SET score=? WHERE id_user=?", 0, user.Id)
}

func GetUserScore(user ConnectedUser) PlayerScore {
	str := fmt.Sprintf("%v", user.Username)
	score := PlayerScore{User: str}
	db := InitDatabase("SQL/database.db")
	data, _ := db.Query(`SELECT score FROM ROOM_USERS WHERE (id_user==?)`, user.Id)
	for data.Next() {
		data.Scan(&score.Score)
	}
	return score
}

func formatLead(str string) PlayerScore {
	result := PlayerScore{}
	stri := ""
	for i := range str {
		if string(str[i]) == ";" {
			result.User = stri
			stri = ""
		} else {
			stri += string(str[i])
		}
	}
	s, _ := strconv.Atoi(stri)
	result.Score = s
	return result
}

func unformatLead(score PlayerScore) string {
	str := fmt.Sprintf("%v", score.User)
	return str + ";" + strconv.Itoa(score.Score)
}

func GetLB(roomId int) LeaderBoard {
	result := LeaderBoard{}
	db := InitDatabase("SQL/database.db")
	data, _ := db.Query(`SELECT first, second, third, fourth, fifth FROM LEADERBOARDS WHERE (room_id==?)`, roomId)
	for data.Next() {
		str := ""
		str1 := ""
		str2 := ""
		str3 := ""
		str4 := ""
		data.Scan(&str, &str1, &str2, &str3, &str4)
		if str != "" {
			fstr := formatLead(str)
			result.UserList = append(result.UserList, fstr)
		}
		if str1 != "" {
			fstr := formatLead(str1)
			result.UserList = append(result.UserList, fstr)
		}
		if str2 != "" {
			fstr := formatLead(str2)
			result.UserList = append(result.UserList, fstr)
		}
		if str3 != "" {
			fstr := formatLead(str3)
			result.UserList = append(result.UserList, fstr)
		}
		if str4 != "" {
			fstr := formatLead(str4)
			result.UserList = append(result.UserList, fstr)
		}
	}
	return result
}

func UptLead(roomId int, user PlayerScore) LeaderBoard {
	actualLB := GetLB(roomId)
	result := LeaderBoard{}
	userAdded := false
	if len(actualLB.UserList) == 0 {
		actualLB.UserList = append(actualLB.UserList, user)
	} else {
		for _, e := range actualLB.UserList {
			fmt.Println(e.User)
			if user.Score > e.Score && user.User != e.User && !userAdded {
				fmt.Println("Test1")
				result.UserList = append(result.UserList, user)
				result.UserList = append(result.UserList, e)
				userAdded = true
			} else if user.Score >= e.Score && user.User == e.User && !userAdded {
				fmt.Println("Test2")
				result.UserList = append(result.UserList, user)
				userAdded = true
			} else if user.User != e.User {
				fmt.Println("Test3")
				result.UserList = append(result.UserList, e)
			} else {
				fmt.Println("Test4")
				result.UserList = append(result.UserList, e)
				userAdded = true
			}
		}
		if !userAdded {
			fmt.Println("Test5")
			result.UserList = append(result.UserList, user)
			userAdded = true
		}
		if len(result.UserList) >= 5 {
			fmt.Println("Test6")
			actualLB.UserList = result.UserList[0:5]
		} else {
			fmt.Println("Test7")
			actualLB.UserList = result.UserList
		}
	}
	return actualLB
}

func SaveLB(roomId int, lb LeaderBoard) {
	db := InitDatabase("SQL/database.db")
	fmt.Println(lb)
	db.Exec("UPDATE LEADERBOARDS SET fifth=?, fourth=?, third=?, second=?, first=? WHERE room_id=?", "", "", "", "", "", roomId)
	for i := range lb.UserList {
		if i == 0 {
			db.Exec("UPDATE LEADERBOARDS SET first=? WHERE room_id=?", unformatLead(lb.UserList[i]), roomId)
		} else if i == 1 {
			db.Exec("UPDATE LEADERBOARDS SET second=? WHERE room_id=?", unformatLead(lb.UserList[i]), roomId)
		} else if i == 2 {
			db.Exec("UPDATE LEADERBOARDS SET third=? WHERE room_id=?", unformatLead(lb.UserList[i]), roomId)
		} else if i == 3 {
			db.Exec("UPDATE LEADERBOARDS SET fourth=? WHERE room_id=?", unformatLead(lb.UserList[i]), roomId)
		} else if i == 4 {
			db.Exec("UPDATE LEADERBOARDS SET fifth=? WHERE room_id=?", unformatLead(lb.UserList[i]), roomId)
		}
	}
}

// Region End - Game Modes

// Region Start - User

func IsAlreadyPlaying(db *sql.DB, user ConnectedUser) bool {
	data, _ := db.Query("SELECT id_user FROM ROOM_USERS WHERE id_user=?", user.Id)
	n := 0
	for data.Next() {
		n++
	}
	return n > 0
}

func ModifyPass(db *sql.DB, userId int, newPass string) string {
	err := checkPass(newPass)
	return err
}

func ModifyUsername(db *sql.DB, userId int, newName string) string {
	err := checkUsername(db, newName)
	return err
}

func ModifyMail(db *sql.DB, userId int, newMail string) string {
	err := checkMail(db, newMail)
	return err
}

func LoginIn(db *sql.DB, prompt string, password string) (ConnectedUser, string) {
	var user ConnectedUser
	var errBis ConnectedUser
	var pass string
	var data *sql.Rows
	var err error
	ret, _ := verifier.Verify(prompt)
	if ret.Syntax.Valid {
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

func LogOut() ConnectedUser {
	var user ConnectedUser
	return user
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

func GetUserById_interface(db *sql.DB, id interface{}) ConnectedUser {
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
		return "Invalid username, length must be greater than 3."
	}
	ret, _ := verifier.Verify(Username)
	if ret.Syntax.Valid {
		return "The username can't be a mail address."
	}
	return ""
}

func checkMail(db *sql.DB, email string) string {
	ret, _ := verifier.Verify(email)
	if !ret.Syntax.Valid {
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
