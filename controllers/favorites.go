package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"projet-groupie/models"
	"projet-groupie/utils"
)

// FavoritesController gère la page des favoris
func FavoritesController(w http.ResponseWriter, r *http.Request) {
	// Récupérer la liste des favoris
	favorites := models.GetFavorites()
	favoriteIDs := favorites.GetAll()

	// Récupérer les cartes correspondantes
	favoriteCards, err := models.GetCardsByIDs(favoriteIDs)
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Erreur Serveur", "Impossible de récupérer les cartes favorites")
		return
	}

	// Marquer toutes les cartes comme favorites (puisqu'elles sont dans la liste des favoris)
	for i := range favoriteCards {
		favoriteCards[i].IsFavorite = true
	}

	// Données pour le template
	data := struct {
		CurrentPage string
		Cards       []models.Card
	}{
		CurrentPage: "favoris",
		Cards:       favoriteCards,
	}

	// Charger et exécuter le template
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		utils.HandleError(err, "Erreur lors du chargement des templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "favorites", data); err != nil {
		utils.HandleError(err, "Erreur lors du rendu du template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}

// Structure pour la requête de modification des favoris
type ToggleFavoriteRequest struct {
	CardID string `json:"cardId"`
}

// Structure pour la réponse de modification des favoris
type ToggleFavoriteResponse struct {
	Success    bool   `json:"success"`
	IsFavorite bool   `json:"isFavorite"`
	Message    string `json:"message,omitempty"`
}

// ToggleFavoriteController gère l'ajout/retrait d'une carte des favoris
func ToggleFavoriteController(w http.ResponseWriter, r *http.Request) {
	// Vérifier la méthode HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Décoder la requête JSON
	var req ToggleFavoriteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ToggleFavoriteResponse{
			Success: false,
			Message: "Format de requête invalide",
		})
		return
	}

	// Vérifier que l'ID de la carte est valide
	if req.CardID == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ToggleFavoriteResponse{
			Success: false,
			Message: "ID de carte non spécifié",
		})
		return
	}

	// Récupérer la liste des favoris
	favorites := models.GetFavorites()

	// Ajouter/retirer la carte des favoris
	isFavorite, err := favorites.Toggle(req.CardID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ToggleFavoriteResponse{
			Success: false,
			Message: "Erreur lors de la modification des favoris",
		})
		return
	}

	// Répondre avec succès
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToggleFavoriteResponse{
		Success:    true,
		IsFavorite: isFavorite,
	})
}

// Structure pour la réponse de suppression des favoris
type ClearFavoritesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ClearFavoritesController gère la suppression de tous les favoris
func ClearFavoritesController(w http.ResponseWriter, r *http.Request) {
	// Vérifier la méthode HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer la liste des favoris
	favorites := models.GetFavorites()

	// Supprimer tous les favoris
	err := favorites.Clear()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ClearFavoritesResponse{
			Success: false,
			Message: "Erreur lors de la suppression des favoris",
		})
		return
	}

	// Répondre avec succès
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ClearFavoritesResponse{
		Success: true,
	})
}
