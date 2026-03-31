package controllers

import (
	"html/template"
	"log"
	"sync"
)

var (
	tmplOnce      sync.Once
	cachedTemplates *template.Template
	tmplErr       error
)

func getTemplates() (*template.Template, error) {
	tmplOnce.Do(func() {
		cachedTemplates, tmplErr = template.ParseGlob("templates/*.html")
		if tmplErr != nil {
			log.Printf("Erreur chargement templates: %v", tmplErr)
		}
	})
	return cachedTemplates, tmplErr
}
