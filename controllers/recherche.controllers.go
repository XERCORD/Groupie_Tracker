package controllers

import (
	"encoding/json"
	"net/http"
	"projet-groupie/templates"
	"strconv"
	"strings"
)

// Nouvelle structure pour les cartes
type Card struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Image       string   `json:"image"`
	Rarity      string   `json:"rarity"`
	Category    string   `json:"category"`
	Series      string   `json:"series"`
	Description string   `json:"description"`
	Types       []string `json:"types"`
}

// Handler pour la recherche
func SearchControler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres de recherche
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	pageStr := r.URL.Query().Get("page")
	pageSize := 20

	// Pagination
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Récupérer toutes les cartes
	resp, err := http.Get("https://api.tcgdex.net/v2/fr/cards")
	if err != nil {
		http.Error(w, "API indisponible", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	var cards []Card
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		http.Error(w, "Erreur de décodage", http.StatusInternalServerError)
		return
	}

	// Filtrer les résultats
	filtered := make([]Card, 0)
	for _, card := range cards {
		matchName := strings.Contains(strings.ToLower(card.Name), strings.ToLower(query))
		matchCategory := category == "" || card.Category == category

		if matchName && matchCategory {
			filtered = append(filtered, card)
		}
	}

	// Pagination
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	results := struct {
		Cards        []Card
		Query        string
		Category     string
		Page         int
		TotalPages   int
		PreviousPage int // Ajouté
		NextPage     int // Ajouté
	}{
		filtered[start:end],
		query,
		category,
		page,
		(len(filtered) + pageSize - 1) / pageSize,
		page - 1, // Calcul de la page précédente
		page + 1, // Calcul de la page suivante
	}

	templates.ListTemp.ExecuteTemplate(w, "recherche", results)
}
