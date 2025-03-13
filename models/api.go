package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Constants pour les URLs de l'API
const (
	BaseURL      = "https://api.tcgdex.net/v2/fr"
	CardsURL     = BaseURL + "/cards"
	SeriesURL    = BaseURL + "/series"
	SetsURL      = BaseURL + "/sets"
	CardURL      = BaseURL + "/cards/%s" // %s sera remplacé par l'ID de la carte
	SeriesSetURL = BaseURL + "/sets/%s"  // %s sera remplacé par l'ID de la série
)

// APIClient représente un client pour l'API
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAPIClient crée un nouveau client API
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Client HTTP avec timeout par défaut
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// GetRequest effectue une requête GET vers l'URL spécifiée
func (c *APIClient) GetRequest(endpoint string) ([]byte, error) {
	// Construit l'URL complète
	url := c.BaseURL + endpoint

	// Effectue la requête HTTP
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête API: %w", err)
	}
	defer resp.Body.Close()

	// Vérifie le code de statut
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API a retourné le statut: %d", resp.StatusCode)
	}

	// Lit le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la réponse: %w", err)
	}

	return body, nil
}

// IsAPIAvailable vérifie si l'API est disponible
func IsAPIAvailable(url string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// ParseQueryParams construit une chaîne de requête à partir d'une map de paramètres
func ParseQueryParams(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	query := "?"
	first := true
	for key, value := range params {
		if !first {
			query += "&"
		}
		query += key + "=" + value
		first = false
	}

	return query
}

// GetAllCards récupère toutes les cartes depuis l'API
func GetAllCards() ([]Card, error) {
	// Effectue la requête HTTP
	resp, err := httpClient.Get(CardsURL)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête API: %w", err)
	}
	defer resp.Body.Close()

	// Vérifie le code de statut
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API a retourné le statut: %d", resp.StatusCode)
	}

	// Lit le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la réponse: %w", err)
	}

	// Désérialise le JSON
	var cards []Card
	if err := json.Unmarshal(body, &cards); err != nil {
		return nil, fmt.Errorf("erreur lors de la désérialisation: %w", err)
	}

	return cards, nil
}

// GetCardByID récupère une carte spécifique par son ID
func GetCardByID(cardID string) (Card, error) {
	// Construit l'URL
	url := fmt.Sprintf(CardURL, cardID)

	// Effectue la requête HTTP
	resp, err := httpClient.Get(url)
	if err != nil {
		return Card{}, fmt.Errorf("erreur lors de la requête API: %w", err)
	}
	defer resp.Body.Close()

	// Vérifie le code de statut
	if resp.StatusCode != http.StatusOK {
		return Card{}, fmt.Errorf("API a retourné le statut: %d", resp.StatusCode)
	}

	// Lit le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Card{}, fmt.Errorf("erreur lors de la lecture de la réponse: %w", err)
	}

	// Désérialise le JSON
	var card Card
	if err := json.Unmarshal(body, &card); err != nil {
		return Card{}, fmt.Errorf("erreur lors de la désérialisation: %w", err)
	}

	return card, nil
}

// GetAllSeries récupère toutes les séries depuis l'API
func GetAllSeries() ([]Series, error) {
	// Effectue la requête HTTP
	resp, err := httpClient.Get(SeriesURL)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête API: %w", err)
	}
	defer resp.Body.Close()

	// Vérifie le code de statut
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API a retourné le statut: %d", resp.StatusCode)
	}

	// Lit le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la réponse: %w", err)
	}

	// Désérialise le JSON
	var series []Series
	if err := json.Unmarshal(body, &series); err != nil {
		return nil, fmt.Errorf("erreur lors de la désérialisation: %w", err)
	}

	return series, nil
}

// GetRecentCards récupère les cartes les plus récentes
func GetRecentCards(limit int) ([]Card, error) {
	// Récupérer toutes les cartes
	allCards, err := GetAllCards()
	if err != nil {
		return nil, err
	}

	// Trier par date de sortie (si disponible)
	// Dans ce cas, nous utilisons simplement les premières cartes
	// car l'API TCGDEX renvoie souvent les plus récentes en premier

	// Si le nombre de cartes est inférieur à la limite, retourne toutes les cartes
	if len(allCards) <= limit {
		return allCards, nil
	}

	// Sinon, retourne les X premières cartes
	return allCards[:limit], nil
}

// GetSeriesByID récupère une série spécifique par son ID
func GetSeriesByID(seriesID string) (Series, error) {
	// Construit l'URL
	url := fmt.Sprintf(SeriesSetURL, seriesID)

	// Effectue la requête HTTP
	resp, err := httpClient.Get(url)
	if err != nil {
		return Series{}, fmt.Errorf("erreur lors de la requête API: %w", err)
	}
	defer resp.Body.Close()

	// Vérifie le code de statut
	if resp.StatusCode != http.StatusOK {
		return Series{}, fmt.Errorf("API a retourné le statut: %d", resp.StatusCode)
	}

	// Lit le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Series{}, fmt.Errorf("erreur lors de la lecture de la réponse: %w", err)
	}

	// Désérialise le JSON
	var series Series
	if err := json.Unmarshal(body, &series); err != nil {
		return Series{}, fmt.Errorf("erreur lors de la désérialisation: %w", err)
	}

	return series, nil
}

// GetCardsBySeries récupère toutes les cartes d'une série spécifique
func GetCardsBySeries(seriesID string) ([]Card, error) {
	// Récupère toutes les cartes
	allCards, err := GetAllCards()
	if err != nil {
		return nil, err
	}

	// Filtre les cartes par série
	var seriesCards []Card
	for _, card := range allCards {
		if card.SeriesID == seriesID {
			seriesCards = append(seriesCards, card)
		}
	}

	return seriesCards, nil
}

// GetCardsByIDs récupère plusieurs cartes par leurs IDs
func GetCardsByIDs(cardIDs []string) ([]Card, error) {
	if len(cardIDs) == 0 {
		return []Card{}, nil
	}

	// Récupère toutes les cartes
	allCards, err := GetAllCards()
	if err != nil {
		return nil, err
	}

	// Crée une map pour accès rapide
	cardMap := make(map[string]Card)
	for _, card := range allCards {
		cardMap[card.ID] = card
	}

	// Récupère les cartes demandées
	var cards []Card
	for _, id := range cardIDs {
		if card, ok := cardMap[id]; ok {
			cards = append(cards, card)
		}
	}

	return cards, nil
}

// SearchCards recherche des cartes selon les critères spécifiés
func SearchCards(query SearchQuery) ([]Card, int, error) {
	// Récupérer toutes les cartes
	allCards, err := GetAllCards()
	if err != nil {
		return nil, 0, err
	}

	// Filtrer les cartes
	var filtered []Card
	for _, card := range allCards {
		// Vérifier si la carte correspond aux critères
		matchName := true
		if query.Query != "" {
			matchName = Contains(ToLowerCase(card.Name), ToLowerCase(query.Query))
		}

		matchCategory := true
		if query.Category != "" {
			matchCategory = card.Category == query.Category
		}

		matchType := true
		if query.Type != "" {
			matchType = false
			for _, t := range card.Types {
				if t == query.Type {
					matchType = true
					break
				}
			}
		}

		matchRarity := true
		if query.Rarity != "" {
			matchRarity = card.Rarity == query.Rarity
		}

		// Si tous les critères correspondent, ajouter la carte
		if matchName && matchCategory && matchType && matchRarity {
			filtered = append(filtered, card)
		}
	}

	// Nombre total de résultats
	totalCount := len(filtered)

	// Pagination
	start := (query.Page - 1) * query.PageSize
	end := start + query.PageSize

	// Vérifier les limites
	if start >= totalCount {
		return []Card{}, totalCount, nil
	}
	if end > totalCount {
		end = totalCount
	}

	return filtered[start:end], totalCount, nil
}

// GetFilterOptions récupère les options de filtres disponibles
func GetFilterOptions() ([]string, []string, []string, error) {
	// Récupère toutes les cartes
	cards, err := GetAllCards()
	if err != nil {
		return nil, nil, nil, err
	}

	// Maps pour stocker les valeurs uniques
	categories := make(map[string]bool)
	types := make(map[string]bool)
	rarities := make(map[string]bool)

	// Parcourt toutes les cartes
	for _, card := range cards {
		// Catégories
		if card.Category != "" {
			categories[card.Category] = true
		}

		// Types
		for _, t := range card.Types {
			if t != "" {
				types[t] = true
			}
		}

		// Raretés
		if card.Rarity != "" {
			rarities[card.Rarity] = true
		}
	}

	// Convertit les maps en slices
	categoryList := make([]string, 0, len(categories))
	for category := range categories {
		categoryList = append(categoryList, category)
	}

	typeList := make([]string, 0, len(types))
	for t := range types {
		typeList = append(typeList, t)
	}

	rarityList := make([]string, 0, len(rarities))
	for rarity := range rarities {
		rarityList = append(rarityList, rarity)
	}

	return categoryList, typeList, rarityList, nil
}

// GetSimilarCards récupère des cartes similaires à une carte donnée
func GetSimilarCards(card Card, limit int) ([]Card, error) {
	// Récupère toutes les cartes
	allCards, err := GetAllCards()
	if err != nil {
		return nil, err
	}

	// Filtre les cartes par catégorie et type
	var similarCards []Card
	for _, c := range allCards {
		// Ne pas inclure la carte elle-même
		if c.ID == card.ID {
			continue
		}

		// Même catégorie
		if c.Category == card.Category {
			// Si c'est un Pokémon, vérifie aussi le type
			if c.Category == "Pokémon" && len(card.Types) > 0 && len(c.Types) > 0 {
				for _, t1 := range card.Types {
					for _, t2 := range c.Types {
						if t1 == t2 {
							similarCards = append(similarCards, c)
							break
						}
					}
				}
			} else {
				similarCards = append(similarCards, c)
			}
		}

		// Si on a atteint la limite, arrête
		if len(similarCards) >= limit {
			break
		}
	}

	return similarCards, nil
}

// Fonctions utilitaires
// ToLowerCase convertit une chaîne en minuscules
func ToLowerCase(s string) string {
	return strings.ToLower(s)
}

// Contains vérifie si une chaîne contient une sous-chaîne
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
