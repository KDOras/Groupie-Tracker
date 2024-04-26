package main

import (
	"database/sql"
	"groupie/src/databaseManager"
	"html/template"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func GamePage(w http.ResponseWriter, r *http.Request) {

	template, err := template.ParseFiles("./gamepage.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./createAccount.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, r)
}

func Login(w http.ResponseWriter, r *http.Request) {

	template, err := template.ParseFiles("./login.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func TrySignIn(w http.ResponseWriter, r *http.Request, db *sql.DB, user *databaseManager.ConnectedUser) {
	r.ParseForm()
	err := databaseManager.CreateNewUser(db, databaseManager.User{Username: r.FormValue("username"), Password: r.FormValue("password"), Email: r.FormValue("email")})
	if err == "" {
		userTry, err := databaseManager.LoggingIn(db, r.FormValue("username"), r.FormValue("password"))
		if err == "" {
			user.Id = userTry.Id
			user.Username = userTry.Username
			http.Redirect(w, r, "/Gamepage", http.StatusSeeOther)
		}
	}
}

func TryLogin(w http.ResponseWriter, r *http.Request, db *sql.DB, user *databaseManager.ConnectedUser) {
	r.ParseForm()
	userTry, err := databaseManager.LoggingIn(db, r.FormValue("username"), r.FormValue("password"))
	if err == "" {
		user.Id = userTry.Id
		user.Username = userTry.Username
		http.Redirect(w, r, "/Gamepage", http.StatusSeeOther)
	}
}

func main() {
	user := databaseManager.ConnectedUser{}
	http.HandleFunc("/", Home)
	http.HandleFunc("/Gamepage", GamePage)
	http.HandleFunc("/Register", func(w http.ResponseWriter, r *http.Request) {
		Create(w, r)
	})
	http.HandleFunc("/Login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r)
	})
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		TryLogin(w, r, databaseManager.InitDatabase("SQL/database.db"), &user)
	})
	http.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		TrySignIn(w, r, databaseManager.InitDatabase("SQL/database.db"), &user)
	})
	fs := http.FileServer(http.Dir("./server/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
