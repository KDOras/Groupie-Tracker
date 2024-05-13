package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {

	url := "https://deezerdevs-deezer.p.rapidapi.com/infos"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "bbae74a3dcmsh814fa6cfa0f2b52p1cd3e5jsnbcf84788fa27")
	req.Header.Add("X-RapidAPI-Host", "deezerdevs-deezer.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}