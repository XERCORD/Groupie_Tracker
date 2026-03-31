package main

import (
	"fmt"
	"log"
	"net/http"
	"projet-groupie/controllers"
	"projet-groupie/models"
)

func main() {
	fmt.Println("Démarrage du serveur TCGDEX Explorer...")

	if !models.IsAPIAvailable(models.BaseURL) {
		log.Fatal("L'API TCGDEX n'est pas disponible. Vérifiez votre connexion internet.")
	}

	fmt.Println("API TCGDEX disponible")

	favorites := models.GetFavorites()
	fmt.Printf("Favoris chargés: %d cartes\n", favorites.Count())

	setupRoutes()

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := "8000"
	fmt.Printf("Serveur démarré sur http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Impossible de démarrer le serveur:", err)
	}
}

func setupRoutes() {
	http.HandleFunc("/", controllers.HomeController)
	http.HandleFunc("/accueil", controllers.HomeController)
	http.HandleFunc("/collection", controllers.CollectionController)
	http.HandleFunc("/recherche", controllers.SearchController)
	http.HandleFunc("/carte/", controllers.CardDetailsController)
	http.HandleFunc("/series/", controllers.SeriesController)
	http.HandleFunc("/favoris", controllers.FavoritesController)
	http.HandleFunc("/a-propos", controllers.AboutController)
	http.HandleFunc("/api/favoris/toggle", controllers.ToggleFavoriteController)
	http.HandleFunc("/api/favoris/clear", controllers.ClearFavoritesController)
	http.HandleFunc("/error", controllers.InternalServerErrorHandler)
	http.HandleFunc("/404", controllers.NotFoundHandler)
}
