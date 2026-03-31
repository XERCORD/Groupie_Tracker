package controllers

import (
	"net/http"
	"projet-groupie/models"
	"projet-groupie/utils"
	"strconv"
)

func SearchController(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	cardType := r.URL.Query().Get("type")
	rarity := r.URL.Query().Get("rarity")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 30 {
		pageSize = 30
	}

	cards, totalResults, err := models.SearchCards(models.SearchQuery{
		Query:    query,
		Category: category,
		Type:     cardType,
		Rarity:   rarity,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Erreur de Recherche", "Impossible de récupérer les résultats de recherche")
		return
	}

	categories, types, rarities, err := models.GetFilterOptions()
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Erreur de Filtres", "Impossible de récupérer les options de filtres")
		return
	}

	favorites := models.GetFavorites()
	for i := range cards {
		cards[i].IsFavorite = favorites.Contains(cards[i].ID)
	}

	totalPages := (totalResults + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

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
		"recherche", query, category, cardType, rarity,
		page, pageSize, page - 1, page + 1, totalPages, totalResults,
		cards, categories, types, rarities,
	}

	tmpl, err := getTemplates()
	if err != nil {
		utils.HandleError(err, "Erreur chargement templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "search", data); err != nil {
		utils.HandleError(err, "Erreur rendu template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
