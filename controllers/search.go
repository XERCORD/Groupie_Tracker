package controllers

import (
	"html/template"
	"net/http"
	"projet-groupie/models"
	"projet-groupie/utils"
	"strconv"
)

// SearchController gère la page de recherche de cartes
func SearchController(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres de recherche
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	cardType := r.URL.Query().Get("type")
	rarity := r.URL.Query().Get("rarity")

	// Gérer la pagination
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	// Limiter la taille de la page
	if pageSize > 30 {
		pageSize = 30
	}

	// Créer la requête de recherche
	searchQuery := models.SearchQuery{
		Query:    query,
		Category: category,
		Type:     cardType,
		Rarity:   rarity,
		Page:     page,
		PageSize: pageSize,
	}

	// Récupérer les cartes filtrées
	cards, totalResults, err := models.SearchCards(searchQuery)
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Erreur de Recherche", "Impossible de récupérer les résultats de recherche")
		return
	}

	// Récupérer les options de filtres
	categories, types, rarities, err := models.GetFilterOptions()
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Erreur de Filtres", "Impossible de récupérer les options de filtres")
		return
	}

	// Récupérer les favoris
	favorites := models.GetFavorites()

	// Mettre à jour le statut de favori pour chaque carte
	for i := range cards {
		cards[i].IsFavorite = favorites.Contains(cards[i].ID)
	}

	// Calculer le nombre total de pages
	totalPages := (totalResults + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	// Données pour le template
	data := struct {
		CurrentPage  string
		Query        string
		Category     string
		Type         string
		Rarity       string
		Page         int
		PageSize     int
		PreviousPage int
		NextPage     int
		TotalPages   int
		TotalResults int
		Cards        []models.Card
		Categories   []string
		Types        []string
		Rarities     []string
	}{
		CurrentPage:  "recherche",
		Query:        query,
		Category:     category,
		Type:         cardType,
		Rarity:       rarity,
		Page:         page,
		PageSize:     pageSize,
		PreviousPage: page - 1,
		NextPage:     page + 1,
		TotalPages:   totalPages,
		TotalResults: totalResults,
		Cards:        cards,
		Categories:   categories,
		Types:        types,
		Rarities:     rarities,
	}

	// Charger et exécuter le template
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		utils.HandleError(err, "Erreur lors du chargement des templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "search", data); err != nil {
		utils.HandleError(err, "Erreur lors du rendu du template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
