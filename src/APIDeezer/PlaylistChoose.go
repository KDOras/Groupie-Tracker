package main

import (
	"encoding/json"
	"net/http"
)

// Structure de données pour stocker les données de la chanson
type Track struct {
    Title   string `json:"title"`
    Preview string `json:"preview"`
}

func main() {
    // Définir une route pour récupérer un extrait d'une chanson spécifique depuis l'API Deezer
    http.HandleFunc("/playTrack", func(w http.ResponseWriter, r *http.Request) {
        // Récupérer le titre de la chanson à jouer depuis les paramètres de la requête
        trackTitle := r.URL.Query().Get("title")

        // Faire une requête à l'API Deezer pour obtenir des informations sur la chanson spécifique
        response, err := http.Get("https://api.deezer.com/search?q=" + trackTitle)
        if err != nil {
            http.Error(w, "Erreur lors de la récupération des informations sur la chanson", http.StatusInternalServerError)
            return
        }
        defer response.Body.Close()
// Structure de données pour décoder la réponse JSON
        var searchResult struct {
            Data []Track `json:"data"`
        }

        // Décoder la réponse JSON
        if err := json.NewDecoder(response.Body).Decode(&searchResult); err != nil {
            http.Error(w, "Erreur lors du décodage des données JSON", http.StatusInternalServerError)
            return
        }

        // Vérifier si des résultats ont été trouvés
        if len(searchResult.Data) == 0 {
            http.Error(w, "Aucun résultat trouvé pour la chanson spécifiée", http.StatusNotFound)
            return
        }

        // Récupérer l'extrait de la première chanson trouvée
        previewURL := searchResult.Data[0].Preview

        // Rediriger l'utilisateur vers l'URL de l'extrait de la chanson
        http.Redirect(w, r, previewURL, http.StatusFound)
    })
}