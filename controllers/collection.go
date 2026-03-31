package controllers

import (
	"net/http"
	"projet-groupie/models"
	"projet-groupie/utils"
	"strconv"
)

func CollectionController(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	serieFilter := r.URL.Query().Get("serie")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "name"
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	const pageSize = 12

	allSeries, err := models.GetAllSeries()
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Erreur Serveur", "Impossible de récupérer les séries")
		return
	}

	serieOpts, err := models.GetSerieOptions()
	if err != nil {
		serieOpts = nil
	}

	filtered := models.FilterSeriesBySerieBlock(models.FilterSeries(allSeries, search), serieFilter)
	sortedSeries := models.SortSeries(filtered, sort)
	totalResults := len(sortedSeries)

	totalPages := (totalResults + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > totalResults {
		end = totalResults
	}

	var paginatedSeries []models.Series
	if start < totalResults {
		paginatedSeries = sortedSeries[start:end]
	} else {
		paginatedSeries = []models.Series{}
	}

	data := struct {
		CurrentPage   string
		Series        []models.Series
		SerieOptions  []models.SerieOption
		Search        string
		SerieFilter   string
		Sort          string
		Page          int
		PageSize      int
		PreviousPage  int
		NextPage      int
		TotalPages    int
		TotalResults  int
	}{"collection", paginatedSeries, serieOpts, search, serieFilter, sort, page, pageSize, page - 1, page + 1, totalPages, totalResults}

	tmpl, err := getTemplates()
	if err != nil {
		utils.HandleError(err, "Erreur chargement templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "collection", data); err != nil {
		utils.HandleError(err, "Erreur rendu template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}

func SeriesController(w http.ResponseWriter, r *http.Request) {
	seriesID := r.URL.Path[len("/series/"):]
	if seriesID == "" {
		RenderErrorPage(w, http.StatusBadRequest, "Erreur", "ID de série non spécifié")
		return
	}

	series, err := models.GetSeriesByID(seriesID)
	if err != nil {
		RenderErrorPage(w, http.StatusNotFound, "Série Non Trouvée", "La série demandée n'existe pas")
		return
	}

	allCards, err := models.GetCardsBySeries(seriesID)
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Erreur Serveur", "Impossible de récupérer les cartes de la série")
		return
	}

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

	totalResults := len(allCards)
	totalPages := (totalResults + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > totalResults {
		end = totalResults
	}

	var paginatedCards []models.Card
	if start < totalResults {
		paginatedCards = allCards[start:end]
	} else {
		paginatedCards = []models.Card{}
	}

	favorites := models.GetFavorites()
	for i := range paginatedCards {
		paginatedCards[i].IsFavorite = favorites.Contains(paginatedCards[i].ID)
	}

	data := struct {
		CurrentPage   string
		Series        models.Series
		Cards         []models.Card
		IsSetFavorite bool
		Page          int
		PageSize      int
		PreviousPage  int
		NextPage      int
		TotalPages    int
		TotalResults  int
	}{"collection", series, paginatedCards, favorites.ContainsSet(series.ID), page, pageSize, page - 1, page + 1, totalPages, totalResults}

	tmpl, err := getTemplates()
	if err != nil {
		utils.HandleError(err, "Erreur chargement templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "series", data); err != nil {
		utils.HandleError(err, "Erreur rendu template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
