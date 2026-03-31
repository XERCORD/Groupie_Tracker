package models

type CardSet struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Card struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Image          string   `json:"image"`
	Rarity         string   `json:"rarity"`
	Category       string   `json:"category"`
	Set            CardSet  `json:"set"`
	SeriesID       string   `json:"-"`
	SeriesName     string   `json:"-"`
	Description    string   `json:"description,omitempty"`
	Types          []string `json:"types,omitempty"`
	ReleaseDate    string   `json:"releaseDate,omitempty"`
	HP             int      `json:"hp,omitempty"`
	Artist         string   `json:"illustrator,omitempty"`
	IsFavorite     bool     `json:"-"`
	LocalizedTypes []string `json:"localizedTypes,omitempty"`
}

type SearchQuery struct {
	Query    string
	Category string
	Type     string
	Rarity   string
	Page     int
	PageSize int
}
