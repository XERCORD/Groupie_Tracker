package models

// Series représente une série de cartes
type Series struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ReleaseDate string    `json:"releaseDate"`
	CardCount   CardCount `json:"cardCount"`
	Logo        string    `json:"logo,omitempty"`
}

// CardCount représente le nombre de cartes dans une série
type CardCount struct {
	Total    int `json:"total"`
	Official int `json:"official"`
}

// SeriesSearchQuery représente les paramètres de recherche pour les séries
type SeriesSearchQuery struct {
	Search   string // Terme de recherche pour le nom de la série
	Sort     string // Option de tri (name, date-desc, date-asc, cards-desc, cards-asc)
	Page     int    // Page courante
	PageSize int    // Nombre d'éléments par page
}

// GetSortOptions retourne les options de tri disponibles pour les séries
func GetSortOptions() []map[string]string {
	return []map[string]string{
		{"id": "name", "name": "Nom"},
		{"id": "date-desc", "name": "Date (récent → ancien)"},
		{"id": "date-asc", "name": "Date (ancien → récent)"},
		{"id": "cards-desc", "name": "Nombre de cartes (↓)"},
		{"id": "cards-asc", "name": "Nombre de cartes (↑)"},
	}
}

// SortSeries trie une liste de séries selon l'option de tri spécifiée
func SortSeries(series []Series, sort string) []Series {
	result := make([]Series, len(series))
	copy(result, series)

	switch sort {
	case "name":
		// Tri par nom (ascendant)
		for i := 0; i < len(result); i++ {
			for j := i + 1; j < len(result); j++ {
				if result[i].Name > result[j].Name {
					result[i], result[j] = result[j], result[i]
				}
			}
		}
	case "date-desc":
		// Tri par date (descendant)
		for i := 0; i < len(result); i++ {
			for j := i + 1; j < len(result); j++ {
				if result[i].ReleaseDate < result[j].ReleaseDate {
					result[i], result[j] = result[j], result[i]
				}
			}
		}
	case "date-asc":
		// Tri par date (ascendant)
		for i := 0; i < len(result); i++ {
			for j := i + 1; j < len(result); j++ {
				if result[i].ReleaseDate > result[j].ReleaseDate {
					result[i], result[j] = result[j], result[i]
				}
			}
		}
	case "cards-desc":
		// Tri par nombre de cartes (descendant)
		for i := 0; i < len(result); i++ {
			for j := i + 1; j < len(result); j++ {
				if result[i].CardCount.Total < result[j].CardCount.Total {
					result[i], result[j] = result[j], result[i]
				}
			}
		}
	case "cards-asc":
		// Tri par nombre de cartes (ascendant)
		for i := 0; i < len(result); i++ {
			for j := i + 1; j < len(result); j++ {
				if result[i].CardCount.Total > result[j].CardCount.Total {
					result[i], result[j] = result[j], result[i]
				}
			}
		}
	default:
		// Par défaut, tri par nom
		for i := 0; i < len(result); i++ {
			for j := i + 1; j < len(result); j++ {
				if result[i].Name > result[j].Name {
					result[i], result[j] = result[j], result[i]
				}
			}
		}
	}

	return result
}

// FilterSeries filtre une liste de séries selon le terme de recherche
func FilterSeries(series []Series, search string) []Series {
	if search == "" {
		return series
	}

	var result []Series
	searchLower := ToLowerCase(search)

	for _, s := range series {
		nameLower := ToLowerCase(s.Name)
		if Contains(nameLower, searchLower) {
			result = append(result, s)
		}
	}

	return result
}
