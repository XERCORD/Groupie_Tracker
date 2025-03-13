package controllers

import (
	"net/http"
	"projet-groupie/templates"
)

// créer une fonction controler qui permet de compléter la fonction route associé au meme template.
func HomeControler(w http.ResponseWriter, r *http.Request) {
	templates.ListTemp.ExecuteTemplate(w, "accueil", nil)
}
