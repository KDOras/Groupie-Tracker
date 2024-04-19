package main

import (
	"groupie/Go/databaseManager"
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

	template, err := template.ParseFiles("./login.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/Gamepage", GamePage)
	http.HandleFunc("/Register", func(w http.ResponseWriter, r *http.Request) {
		Create(w, r, databaseManager.User{})
	})
	http.HandleFunc("/Login", func(w http.ResponseWriter, r *http.Request) {
		Create(w, r, databaseManager.User{})
	})
	fs := http.FileServer(http.Dir("./server/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
