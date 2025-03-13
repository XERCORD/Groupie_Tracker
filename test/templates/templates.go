package templates

import (
	"fmt"
	"html/template"
	"os"
)

// Variable globale qui sert à stocker les templates chargés
var ListTemp *template.Template

// Méthode permettant de charger l'ensemble des templates
func Init() {
	listTemp, tempErr := template.ParseGlob("./templates/*.html")
	if tempErr != nil {
		fmt.Printf("Erreur Template - Une erreur lors du chargement des template \n message d'erreur : %v\n", tempErr.Error())
		os.Exit(1)
	}
	ListTemp = listTemp
}
