package controllers

import (
	"html/template"
	"net/http"
	"projet-groupie/utils"
)

func AboutController(w http.ResponseWriter, r *http.Request) {
	data := struct{ CurrentPage string }{"a-propos"}

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		utils.HandleError(err, "Erreur chargement templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "about", data); err != nil {
		utils.HandleError(err, "Erreur rendu template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
