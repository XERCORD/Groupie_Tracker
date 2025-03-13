package controllers

import (
	"html/template"
	"net/http"
	"projet-groupie/utils"
)

// RenderErrorPage affiche une page d'erreur personnalisée (rendu public avec majuscule)
func RenderErrorPage(w http.ResponseWriter, statusCode int, errorTitle, errorMessage string) {
	// Définir le code de statut HTTP
	w.WriteHeader(statusCode)

	// Données pour le template
	data := struct {
		CurrentPage  string
		StatusCode   int
		ErrorTitle   string
		ErrorMessage string
	}{
		CurrentPage:  "error",
		StatusCode:   statusCode,
		ErrorTitle:   errorTitle,
		ErrorMessage: errorMessage,
	}

	// Charger et exécuter le template
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		// Si le template d'erreur n'est pas trouvé, afficher une erreur simple
		http.Error(w, errorMessage, statusCode)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "error", data); err != nil {
		utils.HandleError(err, "Erreur lors du rendu du template d'erreur")
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}

// NotFoundHandler gère les requêtes pour des pages qui n'existent pas
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	RenderErrorPage(w, http.StatusNotFound, "Page Non Trouvée", "La page que vous recherchez n'existe pas.")
}

// MethodNotAllowedHandler gère les requêtes avec des méthodes non autorisées
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	RenderErrorPage(w, http.StatusMethodNotAllowed, "Méthode Non Autorisée", "La méthode que vous utilisez n'est pas autorisée pour cette ressource.")
}

// InternalServerErrorHandler gère les erreurs internes du serveur
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	RenderErrorPage(w, http.StatusInternalServerError, "Erreur Serveur", "Une erreur interne est survenue sur le serveur.")
}
