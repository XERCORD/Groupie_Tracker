package models

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Initialisation du générateur de nombres aléatoires
func init() {
	rand.Seed(time.Now().UnixNano())
}

// TruncateString tronque une chaîne à la longueur spécifiée et ajoute "..." si nécessaire
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// FormatDate formate une date au format spécifié
func FormatDate(dateStr string) string {
	// Analyse la date (format supposé: "YYYY-MM-DD")
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		// Si l'analyse échoue, retourne la chaîne d'origine
		return dateStr
	}

	// Formate la date
	return date.Format("02/01/2006")
}

// GetRandomItems retourne n éléments aléatoires d'une slice
func GetRandomItems[T any](items []T, n int) []T {
	// Si n est supérieur à la longueur de la slice, retourne tous les éléments
	if n >= len(items) {
		result := make([]T, len(items))
		copy(result, items)
		return result
	}

	// Crée une copie pour éviter de modifier l'original
	temp := make([]T, len(items))
	copy(temp, items)

	// Mélange la slice
	rand.Shuffle(len(temp), func(i, j int) {
		temp[i], temp[j] = temp[j], temp[i]
	})

	// Retourne les n premiers éléments
	return temp[:n]
}

// HandleError gère une erreur en l'enregistrant et en retournant un message
func HandleError(err error, defaultMsg string) string {
	if err != nil {
		log.Printf("Erreur: %v", err)
		return fmt.Sprintf("%s: %v", defaultMsg, err)
	}
	return ""
}

// IsStringInSlice vérifie si une chaîne est présente dans une slice
func IsStringInSlice(value string, slice []string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// GetPageRange génère une plage de numéros de page autour de la page courante
func GetPageRange(currentPage, totalPages, maxVisible int) []int {
	// Si le nombre total de pages est inférieur ou égal au maximum visible,
	// retourne toutes les pages
	if totalPages <= maxVisible {
		pages := make([]int, totalPages)
		for i := 0; i < totalPages; i++ {
			pages[i] = i + 1
		}
		return pages
	}

	// Calcul du nombre de pages à afficher de chaque côté
	sidePages := (maxVisible - 1) / 2

	// Calcul des bornes
	start := currentPage - sidePages
	if start < 1 {
		start = 1
	}

	end := start + maxVisible - 1
	if end > totalPages {
		end = totalPages
		start = end - maxVisible + 1
		if start < 1 {
			start = 1
		}
	}

	// Création de la plage
	pages := make([]int, end-start+1)
	for i := start; i <= end; i++ {
		pages[i-start] = i
	}

	return pages
}
