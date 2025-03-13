package controllers

import (
	"html/template"
	"net/http"
	"projet-groupie/models"
	"projet-groupie/utils"
	"strconv"
)

// CollectionController gère la page de collection (séries)
func CollectionController(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres de filtrage
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "name" // Tri par défaut
	}

	// Récupérer les paramètres de pagination
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize := 12 // Nombre de séries par page

	// Récupérer toutes les séries
	allSeries, err := models.GetAllSeries()
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Erreur Serveur", "Impossible de récupérer les séries")
		return
	}

	// Filtrer les séries selon le terme de recherche
	filteredSeries := models.FilterSeries(allSeries, search)

	// Trier les séries selon l'option de tri
	sortedSeries := models.SortSeries(filteredSeries, sort)

	// Nombre total de résultats
	totalResults := len(sortedSeries)

	// Calculer le nombre total de pages
	totalPages := (totalResults + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	// Vérifier si la page demandée est valide
	if page > totalPages {
		page = totalPages
	}

	// Pagination
	var paginatedSeries []models.Series
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= totalResults {
		paginatedSeries = []models.Series{}
	} else {
		if end > totalResults {
			end = totalResults
		}
		paginatedSeries = sortedSeries[start:end]
	}

	// Données pour le template
	data := struct {
		CurrentPage  string
		Series       []models.Series
		Search       string
		Sort         string
		Page         int
		PageSize     int
		PreviousPage int
		NextPage     int
		TotalPages   int
		TotalResults int
	}{
		CurrentPage:  "collection",
		Series:       paginatedSeries,
		Search:       search,
		Sort:         sort,
		Page:         page,
		PageSize:     pageSize,
		PreviousPage: page - 1,
		NextPage:     page + 1,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}

	// Charger et exécuter le template
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		utils.HandleError(err, "Erreur lors du chargement des templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "collection", data); err != nil {
		utils.HandleError(err, "Erreur lors du rendu du template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}

// SeriesController gère la page de détails d'une série
func SeriesController(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de la série à partir de l'URL
	seriesID := r.URL.Path[len("/series/"):]
	if seriesID == "" {
		RenderErrorPage(w, http.StatusBadRequest, "Erreur", "ID de série non spécifié")
		return
	}

	// Récupérer les informations sur la série
	series, err := models.GetSeriesByID(seriesID)
	if err != nil {
		RenderErrorPage(w, http.StatusNotFound, "Série Non Trouvée", "La série demandée n'existe pas")
		return
	}

	// Récupérer les cartes de la série
	allCards, err := models.GetCardsBySeries(seriesID)
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Erreur Serveur", "Impossible de récupérer les cartes de la série")
		return
	}

	// Récupérer les paramètres de pagination
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

	// Nombre total de résultats
	totalResults := len(allCards)

	// Calculer le nombre total de pages
	totalPages := (totalResults + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	// Vérifier si la page demandée est valide
	if page > totalPages {
		page = totalPages
	}

	// Pagination
	var paginatedCards []models.Card
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= totalResults {
		paginatedCards = []models.Card{}
	} else {
		if end > totalResults {
			end = totalResults
		}
		paginatedCards = allCards[start:end]
	}

	// Vérifier les favoris pour chaque carte
	favorites := models.GetFavorites()
	for i := range paginatedCards {
		paginatedCards[i].IsFavorite = favorites.Contains(paginatedCards[i].ID)
	}

	// Données pour le template
	data := struct {
		CurrentPage  string
		Series       models.Series
		Cards        []models.Card
		Page         int
		PageSize     int
		PreviousPage int
		NextPage     int
		TotalPages   int
		TotalResults int
	}{
		CurrentPage:  "collection",
		Series:       series,
		Cards:        paginatedCards,
		Page:         page,
		PageSize:     pageSize,
		PreviousPage: page - 1,
		NextPage:     page + 1,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}

	// Charger et exécuter le template
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		utils.HandleError(err, "Erreur lors du chargement des templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "series", data); err != nil {
		utils.HandleError(err, "Erreur lors du rendu du template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
