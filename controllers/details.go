package controllers

import (
	"html/template"
	"net/http"
	"projet-groupie/models"
	"projet-groupie/utils"
)

// CardDetailsController gère la page de détails d'une carte
func CardDetailsController(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de la carte à partir de l'URL
	cardID := r.URL.Path[len("/carte/"):]
	if cardID == "" {
		RenderErrorPage(w, http.StatusBadRequest, "Erreur", "ID de carte non spécifié")
		return
	}

	// Récupérer la carte
	card, err := models.GetCardByID(cardID)
	if err != nil {
		RenderErrorPage(w, http.StatusNotFound, "Carte Non Trouvée", "La carte demandée n'existe pas")
		return
	}

	// Vérifier si la carte est dans les favoris
	favorites := models.GetFavorites()
	card.IsFavorite = favorites.Contains(card.ID)

	// Récupérer des cartes similaires
	similarCards, err := models.GetSimilarCards(card, 6)
	if err != nil {
		// Ne pas afficher d'erreur, simplement ne pas montrer de cartes similaires
		similarCards = []models.Card{}
	}

	// Vérifier les favoris pour chaque carte similaire
	for i := range similarCards {
		similarCards[i].IsFavorite = favorites.Contains(similarCards[i].ID)
	}

	// Données pour le template
	data := struct {
		CurrentPage  string
		Card         models.Card
		SimilarCards []models.Card
	}{
		CurrentPage:  "details",
		Card:         card,
		SimilarCards: similarCards,
	}

	// Charger et exécuter le template
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		utils.HandleError(err, "Erreur lors du chargement des templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "details", data); err != nil {
		utils.HandleError(err, "Erreur lors du rendu du template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
