package models

// Card représente une carte Pokémon
type Card struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Image          string   `json:"image"`
	Rarity         string   `json:"rarity"`
	Category       string   `json:"category"`
	SeriesID       string   `json:"setId"`
	SeriesName     string   `json:"set,omitempty"`
	Description    string   `json:"description,omitempty"`
	Types          []string `json:"types,omitempty"`
	ReleaseDate    string   `json:"releaseDate,omitempty"`
	HP             int      `json:"hp,omitempty"`
	Artist         string   `json:"artist,omitempty"`
	IsFavorite     bool     `json:"-"` // Champ non sérialisé pour l'API
	LocalizedTypes []string `json:"localizedTypes,omitempty"`
}

// SearchQuery représente les paramètres de recherche
type SearchQuery struct {
	Query    string // Terme de recherche
	Category string // Catégorie (Pokémon, Énergie, Dresseur)
	Type     string // Type de Pokémon (Feu, Eau, etc.)
	Rarity   string // Rareté de la carte
	Page     int    // Page courante
	PageSize int    // Nombre d'éléments par page
}

// SortOption représente une option de tri pour les cartes
type SortOption struct {
	ID    string
	Name  string
	Order string // asc ou desc
}

// CardAPIParams représente les paramètres pour les requêtes à l'API
type CardAPIParams struct {
	SearchQuery
	SeriesID string   // ID de la série
	IDs      []string // Liste d'IDs pour récupérer des cartes spécifiques
	Sort     string   // Option de tri
}

// PaginatedResult représente un résultat paginé
type PaginatedResult struct {
	Items        interface{} // Items (cartes, séries, etc.)
	TotalCount   int         // Nombre total d'éléments
	Page         int         // Page courante
	PageSize     int         // Nombre d'éléments par page
	TotalPages   int         // Nombre total de pages
	NextPage     int         // Page suivante
	PreviousPage int         // Page précédente
}

// NewPaginatedResult crée un nouveau résultat paginé
func NewPaginatedResult(items interface{}, totalCount, page, pageSize int) PaginatedResult {
	// Calcul du nombre total de pages
	totalPages := (totalCount + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	// Calcul des pages précédente et suivante
	previousPage := page - 1
	if previousPage < 1 {
		previousPage = 1
	}

	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	return PaginatedResult{
		Items:        items,
		TotalCount:   totalCount,
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
		NextPage:     nextPage,
		PreviousPage: previousPage,
	}
}
