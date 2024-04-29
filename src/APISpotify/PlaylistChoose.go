package PlaylistChoose

import (
	"fmt"

	"github.com/zmb3/spotify"
)

func ChoosePlaylist(client *spotify.Client) {
	// Get the current user
	playlists, err := client.Search("2024", spotify.SearchTypePlaylist)
	if err != nil {
		panic(err)
	}

	// Print the user's playlists
	for _, playlist := range playlists.Playlists.Playlists {
		fmt.Println(playlist.Name)
	}
}
