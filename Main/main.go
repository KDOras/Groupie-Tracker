package main

import (
	"groupie/Go/accountManager"
	"html/template"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./html/index.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func main() {
	accountManager.LoggingIn("test", "test")
	print("test1")
	http.HandleFunc("/", Home)

	fs := http.FileServer(http.Dir("./Main/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
