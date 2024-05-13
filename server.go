package main

import (
	"database/sql"
	"groupie/src/databaseManager"
	Game "groupie/src/games"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

var tpl *template.Template

var store = sessions.NewCookieStore([]byte("super-secret-password"))

type Err struct {
	Err string
}

type GamePageVar struct {
	IsSidePanelOpen interface{}
	Username        interface{}
	RoomList        []databaseManager.Room
}

type DeafVar struct {
	IsStarted   bool
	HasWon      bool
	Song        Game.Song
	Leaderboard databaseManager.LeaderBoard
}

type Category struct {
	Id   int
	Name string
}

type scattergorries struct {
	Letter     string
	Categories []Category
}

func Home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func GamePage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	db := databaseManager.InitDatabase("SQL/database.db")
	data, _ := db.Query("SELECT id, max_player, name, password, id_game FROM ROOMS")
	roomList := []databaseManager.Room{}
	for data.Next() {
		var room databaseManager.Room
		data.Scan(&room.Id, &room.MaxPlayer, &room.Name, &room.Password, &room.GameMode.Id)
		room.NumberOfPlayer = databaseManager.GetNumberOfPlayerFromRoom(db, room.Id)
		room.GameMode = databaseManager.GetGame(db, room.GameMode.Id)
		roomList = append(roomList, room)
	}
	pageVar := GamePageVar{Username: session.Values["Username"], IsSidePanelOpen: session.Values["IsSidePanelOpen"], RoomList: roomList}
	template, err := template.ParseFiles("./gamepage.html")
	if err != nil {
		log.Fatal(err)
	}
	session.Save(r, w)
	if session.Values["IsSidePanelOpen"] == true && session.Values["KeepSidePanelOpen"] == false {
		http.Redirect(w, r, "/openProfile", http.StatusSeeOther)
	} else {
		session.Values["KeepSidePanelOpen"] = false
		session.Save(r, w)
		template.Execute(w, pageVar)
	}
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	userId := session.Values["User-Id"]
	username := session.Values["Username"]
	roomId := r.FormValue("JoinButton")
	formatizedRoomId, _ := strconv.Atoi(roomId)
	err := databaseManager.JoinRoom(databaseManager.InitDatabase("SQL/database.db"), databaseManager.ConnectedUser{Id: userId, Username: username}, databaseManager.GetRoom(databaseManager.InitDatabase("SQL/database.db"), formatizedRoomId))
	room := databaseManager.GetRoom(databaseManager.InitDatabase("SQL/database.db"), formatizedRoomId)
	if err == "" {
		if room.GameMode.Id == 0 {
			http.Redirect(w, r, "/BlindTest", http.StatusSeeOther)
		} else if room.GameMode.Id == 1 {
			http.Redirect(w, r, "/DeafTest", http.StatusSeeOther)
		} else if room.GameMode.Id == 2 {
			http.Redirect(w, r, "/ScatterGorries", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/Gamepage", http.StatusSeeOther)
	}
}

func OpenSidePanel(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	if session.Values["IsSidePanelOpen"] == false {
		session.Values["IsSidePanelOpen"] = true
		session.Values["KeepSidePanelOpen"] = true
	} else {
		session.Values["IsSidePanelOpen"] = false
	}
	session.Save(r, w)
	http.Redirect(w, r, "/Gamepage", http.StatusSeeOther)
}

func Create(w http.ResponseWriter, r *http.Request, logErr *Err) {
	template, err := template.ParseFiles("./createAccount.html")
	if err != nil {
		log.Fatal(err)
	}
	if logErr.Err != "" {
		template.Execute(w, logErr)
		*logErr = Err{}
	} else {
		template.Execute(w, r)
	}
}

func Login(w http.ResponseWriter, r *http.Request, logErr *Err) {
	template, err := template.ParseFiles("./login.html")
	if err != nil {
		log.Fatal(err)
	}
	if logErr.Err != "" {
		template.Execute(w, logErr)
		*logErr = Err{}
	} else {
		template.Execute(w, r)
	}
}

func TrySignIn(w http.ResponseWriter, r *http.Request, db *sql.DB, dbErr *Err) {
	user := databaseManager.ConnectedUser{}
	r.ParseForm()
	err := databaseManager.CreateNewUser(db, databaseManager.User{Username: r.FormValue("username"), Password: r.FormValue("password"), Email: r.FormValue("email")})
	if err == "" {
		userTry, err := databaseManager.LoginIn(db, r.FormValue("username"), r.FormValue("password"))
		if err == "" {
			user.Id = userTry.Id
			user.Username = userTry.Username
			session, _ := store.Get(r, "SessionName")
			session.Options.MaxAge = 87600 * 7
			session.Values["User"] = user
			session.Save(r, w)
			http.Redirect(w, r, "/Gamepage", http.StatusSeeOther)
		}
	} else {
		dbErr.Err = err
		http.Redirect(w, r, "/Register", http.StatusSeeOther)
	}
}

func TryLogin(w http.ResponseWriter, r *http.Request, db *sql.DB, dbErr *Err) {
	session, _ := store.Get(r, "session-name")
	session.Options = &sessions.Options{Path: "/", MaxAge: 86400 * 7, HttpOnly: true}
	user := databaseManager.ConnectedUser{}
	r.ParseForm()
	userTry, err := databaseManager.LoginIn(db, r.FormValue("username"), r.FormValue("password"))
	if err == "" {
		user.Id = userTry.Id
		user.Username = userTry.Username
		session.Values["User-Id"] = user.Id
		session.Values["Username"] = user.Username
		err := session.Save(r, w)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/Gamepage", http.StatusSeeOther)
	} else {
		dbErr.Err = err
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
	}
}

func GoDisco(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/Login", http.StatusSeeOther)
}

func BlindTest(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./blindtest.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, r)
}
func BlindTestGame(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./CG_Blindtest.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, r)
}

func ScatterGorries(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./scattergorries.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.(http.Flusher).Flush()

	template.Execute(w, r)
}

func ScatterGorriesGame(w http.ResponseWriter, r *http.Request, s scattergorries) {
	template, err := template.ParseFiles("./CG_Scattergories.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, s)
}

func DeafTest(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./Deaftest.html")
	if err != nil {
		log.Fatal(err)
	}

	template.Execute(w, r)
}

func CreateDeafRoom(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	user := session.Values["User-Id"]
	formatMaxP, _ := strconv.Atoi(r.FormValue("maxPlayer"))
	room := databaseManager.Room{Name: r.FormValue("roomName"), MaxPlayer: formatMaxP, CreatedBy: databaseManager.GetUserById_interface(databaseManager.InitDatabase("SQL/database.db"), user), GameMode: databaseManager.GameMode{Id: 1}, Password: ""}
	databaseManager.CreateRoom(databaseManager.InitDatabase("SQL/database.db"), room)
	http.Redirect(w, r, "/DeafTest/Start", http.StatusSeeOther)
}

func DeafTest_Start(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./CG_Deaftest.html")
	if err != nil {
		log.Fatal(err)
	}
	vars := DeafVar{IsStarted: true, Song: Game.GetRandomSong()}

	template.Execute(w, vars)
}

func DeafTest_End(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./CG_Deaftest.html")
	if err != nil {
		log.Fatal(err)
	}
	session, _ := store.Get(r, "session-name")
	user := session.Values["User-Id"]
	fuser := databaseManager.GetUserById_interface(databaseManager.InitDatabase("SQL/database.db"), user)
	vars := DeafVar{IsStarted: false, Song: Game.Song{Name: r.FormValue("submitButton")}}
	if r.FormValue("submitButton") == r.FormValue("input") {
		vars.HasWon = true
		databaseManager.IncreaseUserScore(fuser)
	} else {
		vars.HasWon = false
		room := databaseManager.GetRoomFromUser(databaseManager.InitDatabase("SQL/database.db"), fuser)
		newLB := databaseManager.UptLead(room.Id, databaseManager.GetUserScore(fuser))
		databaseManager.SaveLB(room.Id, newLB)
		defer databaseManager.ResetUserScore(fuser)
	}
	vars.Leaderboard = databaseManager.GetLB(databaseManager.GetRoomFromUser(databaseManager.InitDatabase("SQL/database.db"), fuser).Id)
	template.Execute(w, vars)
}

func Settings(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./settings.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, r)
}

func main() {
	anciantLetters := []string{}
	s := scattergorries{Letter: Game.RandomLetter(anciantLetters), Categories: []Category{{Name: "Animal"}, {Name: "Fruit"}, {Name: "Country"}, {Name: "City"}, {Name: "Object"}, {Name: "Name"}}}
	tpl, _ = template.ParseGlob("*.html")
	dbErr := Err{Err: ""}
	http.HandleFunc("/", Home)
	http.HandleFunc("/Gamepage", GamePage)
	http.HandleFunc("/Settings", Settings)
	http.HandleFunc("/BlindTest", BlindTest)
	http.HandleFunc("/BlindTestGame", BlindTestGame)
	http.HandleFunc("/DeafTest", DeafTest)
	http.HandleFunc("/Create/DeafTest", CreateDeafRoom)
	http.HandleFunc("/DeafTest/Start", DeafTest_Start)
	http.HandleFunc("/DeafTest/End", DeafTest_End)
	http.HandleFunc("/ScatterGorries", ScatterGorries)
	http.HandleFunc("/S", func(w http.ResponseWriter, r *http.Request) {
		ScatterGorriesGame(w, r, s)
	})
	http.HandleFunc("/JoinRoom", JoinRoom)
	http.HandleFunc("/Register", func(w http.ResponseWriter, r *http.Request) {
		Create(w, r, &dbErr)
	})
	http.HandleFunc("/Login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, &dbErr)
	})
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		TryLogin(w, r, databaseManager.InitDatabase("SQL/database.db"), &dbErr)
	})
	http.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		TrySignIn(w, r, databaseManager.InitDatabase("SQL/database.db"), &dbErr)
	})
	http.HandleFunc("/disconnect", func(w http.ResponseWriter, r *http.Request) {
		GoDisco(w, r)
	})
	http.HandleFunc("/openProfile", func(w http.ResponseWriter, r *http.Request) {
		OpenSidePanel(w, r)
	})
	fs := http.FileServer(http.Dir("./server/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
