package controllers

import (
	"html/template"
	"net/http"
	"projet-groupie/models"
	"projet-groupie/utils"
)

// HomeController gère la page d'accueil
func HomeController(w http.ResponseWriter, r *http.Request) {
	// Rediriger / vers /accueil
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	// Récupérer les cartes récentes
	recentCards, err := models.GetRecentCards(6)
	if err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Erreur Serveur", "Impossible de récupérer les cartes récentes")
		return
	}

	// Vérifier les favoris pour chaque carte
	favorites := models.GetFavorites()
	for i := range recentCards {
		recentCards[i].IsFavorite = favorites.Contains(recentCards[i].ID)
	}

	// Données pour le template
	data := struct {
		CurrentPage string
		RecentCards []models.Card
	}{
		CurrentPage: "accueil",
		RecentCards: recentCards,
	}

	// Charger et exécuter le template
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		utils.HandleError(err, "Erreur lors du chargement des templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "home", data); err != nil {
		utils.HandleError(err, "Erreur lors du rendu du template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
