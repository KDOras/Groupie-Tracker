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

func Create(w http.ResponseWriter, r *http.Request, RegisterVar databaseManager.User) {
	template, err := template.ParseFiles("./createAccount.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, RegisterVar)
}

func Login(w http.ResponseWriter, r *http.Request, LoginVar databaseManager.User) {

	template, err := template.ParseFiles("./signin.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Action(w http.ResponseWriter, r *http.Request, db *sql.DB, user *databaseManager.ConnectedUser) {
	r.ParseForm()
	if r.URL.Path == "/Register" {
		user := databaseManager.User{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			Email:    r.FormValue("email"),
		}
		databaseManager.CreateNewUser(db, user)
	} else if r.URL.Path == "/Login" {

	}
}

func TryLogin(w http.ResponseWriter, r *http.Request, db *sql.DB, user *databaseManager.ConnectedUser) {
	r.ParseForm()
	inputUsername := r.FormValue("floatingInput")
	inputPassword := r.FormValue("floatingPassword")
	userTry, err := databaseManager.LoggingIn(db, inputUsername, inputPassword)
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
		Create(w, r, databaseManager.User{})
	})
	http.HandleFunc("/Login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, databaseManager.User{})
	})
	http.HandleFunc("/EAction", func(w http.ResponseWriter, r *http.Request) {
		TryLogin(w, r, databaseManager.InitDatabase("SQL/database.db"), &user)
	})
	fs := http.FileServer(http.Dir("./server/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
