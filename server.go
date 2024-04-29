package main

import (
	"database/sql"
	"fmt"
	"groupie/src/databaseManager"
	"html/template"
	"log"
	"net/http"
)

type Err struct {
	Err string
}

func Home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func GamePage(w http.ResponseWriter, r *http.Request, user *databaseManager.ConnectedUser) {
	template, err := template.ParseFiles("./gamepage.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, user)
}

func Create(w http.ResponseWriter, r *http.Request, logErr *Err) {
	template, err := template.ParseFiles("./createAccount.html")
	if err != nil {
		log.Fatal(err)
	}
	if logErr.Err != "" {
		template.Execute(w, logErr)
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
	} else {
		template.Execute(w, r)
	}
}

func TrySignIn(w http.ResponseWriter, r *http.Request, db *sql.DB, user *databaseManager.ConnectedUser, dbErr *Err) {
	r.ParseForm()
	err := databaseManager.CreateNewUser(db, databaseManager.User{Username: r.FormValue("username"), Password: r.FormValue("password"), Email: r.FormValue("email")})
	fmt.Println(r.FormValue("username"))
	if err == "" {
		userTry, err := databaseManager.LoginIn(db, r.FormValue("username"), r.FormValue("password"))
		if err == "" {
			user.Id = userTry.Id
			user.Username = userTry.Username
			http.Redirect(w, r, "/Gamepage", http.StatusSeeOther)
		}
	} else {
		dbErr.Err = err
		http.Redirect(w, r, "/Register", http.StatusSeeOther)
	}
}

func TryLogin(w http.ResponseWriter, r *http.Request, db *sql.DB, user *databaseManager.ConnectedUser, dbErr *Err) {
	r.ParseForm()
	userTry, err := databaseManager.LoginIn(db, r.FormValue("username"), r.FormValue("password"))
	if err == "" {
		user.Id = userTry.Id
		user.Username = userTry.Username
		http.Redirect(w, r, "/Gamepage", http.StatusSeeOther)
	} else {
		dbErr.Err = err
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
	}
}

func main() {
	user := databaseManager.ConnectedUser{}
	dbErr := Err{Err: ""}
	http.HandleFunc("/", Home)
	http.HandleFunc("/Gamepage", func(w http.ResponseWriter, r *http.Request) {
		GamePage(w, r, &user)
	})
	http.HandleFunc("/Register", func(w http.ResponseWriter, r *http.Request) {
		Create(w, r, &dbErr)
	})
	http.HandleFunc("/Login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, &dbErr)
	})
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		TryLogin(w, r, databaseManager.InitDatabase("SQL/database.db"), &user, &dbErr)
	})
	http.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		TrySignIn(w, r, databaseManager.InitDatabase("SQL/database.db"), &user, &dbErr)
	})
	fs := http.FileServer(http.Dir("./server/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
