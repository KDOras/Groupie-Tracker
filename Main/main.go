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
	http.HandleFunc("/", Home)
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./html/"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./server/"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
