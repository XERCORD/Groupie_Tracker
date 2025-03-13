package controllers

import (
	"html/template"
	"net/http"
	"projet-groupie/utils"
)

// AboutController gère la page À propos
func AboutController(w http.ResponseWriter, r *http.Request) {
	// Données pour le template
	data := struct {
		CurrentPage string
	}{
		CurrentPage: "a-propos",
	}

	// Charger et exécuter le template
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		utils.HandleError(err, "Erreur lors du chargement des templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "about", data); err != nil {
		utils.HandleError(err, "Erreur lors du rendu du template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
