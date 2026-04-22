package controllers

import (
	"net/http"
	"net/mail"
	"projet-groupie/models"
	"projet-groupie/utils"
	"strings"
)

const maxContactMessageLen = 8000

func ContactController(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleContactPost(w, r)
		return
	}

	data := struct {
		CurrentPage string
		Success     bool
		Error       string
		Name        string
		Email       string
		Subject     string
		Message     string
	}{
		CurrentPage: "contact",
		Success:     r.URL.Query().Get("sent") == "1",
	}

	tmpl, err := getTemplates()
	if err != nil {
		utils.HandleError(err, "Erreur chargement templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "contact", data); err != nil {
		utils.HandleError(err, "Erreur rendu template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}

func handleContactPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/contact", http.StatusSeeOther)
		return
	}

	name := strings.TrimSpace(r.PostFormValue("name"))
	email := strings.TrimSpace(r.PostFormValue("email"))
	subject := strings.TrimSpace(r.PostFormValue("subject"))
	message := strings.TrimSpace(r.PostFormValue("message"))

	errMsg := validateContact(name, email, subject, message)
	if errMsg != "" {
		renderContactForm(w, name, email, subject, message, errMsg, false)
		return
	}

	if err := models.SaveContactMessage(name, email, subject, message); err != nil {
		utils.HandleError(err, "Erreur enregistrement message contact")
		renderContactForm(w, name, email, subject, message, "Impossible d'enregistrer votre message pour le moment. Réessayez plus tard.", false)
		return
	}

	http.Redirect(w, r, "/contact?sent=1", http.StatusSeeOther)
}

func validateContact(name, email, subject, message string) string {
	if name == "" {
		return "Le nom est obligatoire."
	}
	if len(name) > 120 {
		return "Le nom est trop long."
	}
	if email == "" {
		return "L'adresse e-mail est obligatoire."
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return "Adresse e-mail invalide."
	}
	if subject == "" {
		return "Le sujet est obligatoire."
	}
	if len(subject) > 200 {
		return "Le sujet est trop long."
	}
	if message == "" {
		return "Le message est obligatoire."
	}
	if len(message) > maxContactMessageLen {
		return "Le message est trop long."
	}
	return ""
}

func renderContactForm(w http.ResponseWriter, name, email, subject, message, errMsg string, success bool) {
	data := struct {
		CurrentPage string
		Success     bool
		Error       string
		Name        string
		Email       string
		Subject     string
		Message     string
	}{"contact", success, errMsg, name, email, subject, message}

	tmpl, err := getTemplates()
	if err != nil {
		utils.HandleError(err, "Erreur chargement templates")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "contact", data); err != nil {
		utils.HandleError(err, "Erreur rendu template")
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
