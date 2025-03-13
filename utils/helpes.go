package utils

import (
	"fmt"
	"log"
	"strings"
)

// HandleError gère une erreur en l'enregistrant dans les logs
func HandleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

// FormatDate formate une date au format français
func FormatDate(dateStr string) string {
	// Si la date est au format ISO 8601 (YYYY-MM-DD)
	if len(dateStr) >= 10 && dateStr[4] == '-' && dateStr[7] == '-' {
		year := dateStr[0:4]
		month := dateStr[5:7]
		day := dateStr[8:10]
		return fmt.Sprintf("%s/%s/%s", day, month, year)
	}
	return dateStr
}

// TruncateString tronque une chaîne si elle est trop longue
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// Capitalize met la première lettre d'une chaîne en majuscule
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FormatPrice formate un prix avec le séparateur de milliers
func FormatPrice(price float64) string {
	return fmt.Sprintf("%.2f €", price)
}

// IsValidURL vérifie si une URL est valide
func IsValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// CleanString nettoie une chaîne de caractères pour un usage dans les URL
func CleanString(s string) string {
	// Remplacer les caractères accentués
	replacer := strings.NewReplacer(
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"à", "a", "â", "a", "ä", "a",
		"î", "i", "ï", "i",
		"ô", "o", "ö", "o",
		"ù", "u", "û", "u", "ü", "u",
		"ç", "c",
		" ", "-",
	)
	cleaned := replacer.Replace(strings.ToLower(s))

	// Supprimer les caractères spéciaux
	cleaned = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, cleaned)

	return cleaned
}
