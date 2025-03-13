package main

import (
	"net/http"
	"projet-groupie/controllers"
)

// setupRoutes configure toutes les routes de l'application
func setupRoutes() {
	// Pages principales
	http.HandleFunc("/", controllers.HomeController)
	http.HandleFunc("/accueil", controllers.HomeController)
	http.HandleFunc("/collection", controllers.CollectionController)
	http.HandleFunc("/recherche", controllers.SearchController)
	http.HandleFunc("/carte/", controllers.CardDetailsController)
	http.HandleFunc("/series/", controllers.SeriesController)
	http.HandleFunc("/favoris", controllers.FavoritesController)
	http.HandleFunc("/a-propos", controllers.AboutController)

	// API pour les favoris
	http.HandleFunc("/api/favoris/toggle", controllers.ToggleFavoriteController)
	http.HandleFunc("/api/favoris/clear", controllers.ClearFavoritesController)

	// Gestion des erreurs
	http.HandleFunc("/error", controllers.InternalServerErrorHandler)
	http.HandleFunc("/404", controllers.NotFoundHandler)
}
