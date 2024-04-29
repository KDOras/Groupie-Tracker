package main

import (
	PlaylistChoose "groupie/src/APISpotify"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

func main() {
	auth := spotify.NewAuthenticator("http://localhost:8080/callback", spotify.ScopeUserReadPrivate)
	ch := make(chan *spotify.Client)
	http.HandleFunc(auth.AuthURL(""), func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token("state", r)
		if err != nil {
			log.Fatalf("Couldn't get token: %v", err)
		}
		client := auth.NewClient(token)
		ch <- &client
	})
	http.ListenAndServe(":8080", nil)
	client := <-ch
	PlaylistChoose.ChoosePlaylist(client)
}
