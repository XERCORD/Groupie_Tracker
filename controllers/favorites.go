package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"projet-groupie/models"
	"projet-groupie/utils"
)

func FavoritesController(w http.ResponseWriter, r *http.Request) {
	favorites := models.GetFavorites()
	favoriteCards, err := models.GetCardsByIDs(favorites.GetAll())
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Erreur Serveur", "Impossible de récupérer les cartes favorites")
		return
	}

	for i := range favoriteCards {
		favoriteCards[i].IsFavorite = true
	}

	data := struct {
		CurrentPage string
		Cards       []models.Card
	}{"favoris", favoriteCards}

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		utils.HandleError(err, "Erreur chargement templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "favorites", data); err != nil {
		utils.HandleError(err, "Erreur rendu template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}

type ToggleFavoriteRequest struct {
	CardID string `json:"cardId"`
}

type ToggleFavoriteResponse struct {
	Success    bool   `json:"success"`
	IsFavorite bool   `json:"isFavorite"`
	Message    string `json:"message,omitempty"`
}

func ToggleFavoriteController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req ToggleFavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ToggleFavoriteResponse{Success: false, Message: "Format de requête invalide"})
		return
	}

	if req.CardID == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ToggleFavoriteResponse{Success: false, Message: "ID de carte non spécifié"})
		return
	}

	isFavorite, err := models.GetFavorites().Toggle(req.CardID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ToggleFavoriteResponse{Success: false, Message: "Erreur lors de la modification des favoris"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToggleFavoriteResponse{Success: true, IsFavorite: isFavorite})
}

type ClearFavoritesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func ClearFavoritesController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := models.GetFavorites().Clear(); err != nil {
		json.NewEncoder(w).Encode(ClearFavoritesResponse{Success: false, Message: "Erreur lors de la suppression des favoris"})
		return
	}

	json.NewEncoder(w).Encode(ClearFavoritesResponse{Success: true})
}
