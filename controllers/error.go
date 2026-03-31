package controllers

import (
	"html/template"
	"net/http"
	"projet-groupie/utils"
)

func RenderErrorPage(w http.ResponseWriter, statusCode int, errorTitle, errorMessage string) {
	w.WriteHeader(statusCode)

	data := struct {
		CurrentPage  string
		StatusCode   int
		ErrorTitle   string
		ErrorMessage string
	}{"error", statusCode, errorTitle, errorMessage}

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		http.Error(w, errorMessage, statusCode)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "error", data); err != nil {
		utils.HandleError(err, "Erreur rendu template d'erreur")
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	RenderErrorPage(w, http.StatusNotFound, "Page Non Trouvée", "La page que vous recherchez n'existe pas.")
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	RenderErrorPage(w, http.StatusMethodNotAllowed, "Méthode Non Autorisée", "La méthode utilisée n'est pas autorisée.")
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	RenderErrorPage(w, http.StatusInternalServerError, "Erreur Serveur", "Une erreur interne est survenue.")
}
