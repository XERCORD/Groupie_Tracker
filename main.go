package main

import (
	"fmt"
	"log"
	"net/http"
	"projet-groupie/controllers"
	"projet-groupie/models"
	"projet-groupie/utils"
)

func main() {
	// Afficher le message de démarrage
	fmt.Println("Démarrage du serveur TCGDEX Explorer...")

	// Vérifier la disponibilité de l'API
	if !models.IsAPIAvailable(models.BaseURL) {
		log.Fatal("L'API TCGDEX n'est pas disponible. Vérifiez votre connexion internet ou réessayez plus tard.")
	}

	// Afficher les informations sur l'API
	fmt.Println("API TCGDEX disponible")

	// Initialiser le fichier des favoris
	favorites := models.GetFavorites()
	fmt.Printf("Favoris chargés: %d cartes\n", favorites.Count())

	// Configurer les routes
	setupRoutes()

	// Configurer le serveur de fichiers statiques
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Démarrer le serveur
	port := "8000"
	fmt.Printf("Serveur démarré sur http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		utils.HandleError(err, "Erreur lors du démarrage du serveur")
		log.Fatal("Impossible de démarrer le serveur:", err)
	}
}

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
