package controllers

import (
	"net/http"
	"projet-groupie/models"
	"projet-groupie/utils"
)

func HomeController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	recentCards, err := models.GetRecentCards(6)
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Erreur Serveur", "Impossible de récupérer les cartes récentes")
		return
	}

	favorites := models.GetFavorites()
	for i := range recentCards {
		recentCards[i].IsFavorite = favorites.Contains(recentCards[i].ID)
	}

	data := struct {
		CurrentPage string
		RecentCards []models.Card
	}{"accueil", recentCards}

	tmpl, err := getTemplates()
	if err != nil {
		utils.HandleError(err, "Erreur chargement templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "home", data); err != nil {
		utils.HandleError(err, "Erreur rendu template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
