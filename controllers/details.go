package controllers

import (
	"net/http"
	"projet-groupie/models"
	"projet-groupie/utils"
)

func CardDetailsController(w http.ResponseWriter, r *http.Request) {
	cardID := r.URL.Path[len("/carte/"):]
	if cardID == "" {
		RenderErrorPage(w, http.StatusBadRequest, "Erreur", "ID de carte non spécifié")
		return
	}

	card, err := models.GetCardByID(cardID)
	if err != nil {
		RenderErrorPage(w, http.StatusNotFound, "Carte Non Trouvée", "La carte demandée n'existe pas")
		return
	}

	favorites := models.GetFavorites()
	card.IsFavorite = favorites.Contains(card.ID)

	similarCards, err := models.GetSimilarCards(card, 6)
	if err != nil {
		similarCards = []models.Card{}
	}

	for i := range similarCards {
		similarCards[i].IsFavorite = favorites.Contains(similarCards[i].ID)
	}

	data := struct {
		CurrentPage  string
		Card         models.Card
		SimilarCards []models.Card
	}{"details", card, similarCards}

	tmpl, err := getTemplates()
	if err != nil {
		utils.HandleError(err, "Erreur chargement templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "details", data); err != nil {
		utils.HandleError(err, "Erreur rendu template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
