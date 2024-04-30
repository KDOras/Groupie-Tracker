package main

import (
	"database/sql"
	"groupie/src/databaseManager"
	"html/template"
	"log"
	"net/http"

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
	pageVar := GamePageVar{Username: session.Values["Username"], IsSidePanelOpen: session.Values["IsSidePanelOpen"]}
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
	http.Redirect(w, r, "/Login", http.StatusSeeOther)
}

func main() {
	tpl, _ = template.ParseGlob("*.html")
	dbErr := Err{Err: ""}
	http.HandleFunc("/", Home)
	http.HandleFunc("/Gamepage", GamePage)
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
