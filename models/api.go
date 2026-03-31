package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	BaseURL      = "https://api.tcgdex.net/v2/fr"
	CardsURL     = BaseURL + "/cards"
	SeriesURL    = BaseURL + "/series"
	SetsURL      = BaseURL + "/sets"
	CardURL      = BaseURL + "/cards/%s"
	SeriesSetURL = BaseURL + "/sets/%s"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func IsAPIAvailable(baseURL string) bool {
	client := &http.Client{Timeout: 10 * time.Second}
	fmt.Printf("Vérification de l'API à : %s/cards\n", baseURL)
	resp, err := client.Get(baseURL + "/cards")
	if err != nil {
		fmt.Printf("Erreur connexion API: %v\n", err)
		return false
	}
	defer resp.Body.Close()
	fmt.Printf("Statut de réponse: %d\n", resp.StatusCode)
	return resp.StatusCode == http.StatusOK
}

func fetchJSON(url string, dest interface{}) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("erreur requête: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statut API: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erreur lecture réponse: %w", err)
	}

	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("erreur désérialisation: %w", err)
	}
	return nil
}

func appendImageSuffix(url, suffix string) string {
	if url == "" || strings.HasSuffix(url, ".webp") || strings.HasSuffix(url, ".png") {
		return url
	}
	return url + suffix
}

func FetchCardDetails(cards []Card) []Card {
	type result struct {
		index int
		card  Card
	}

	results := make(chan result, len(cards))
	sem := make(chan struct{}, 10)

	var wg sync.WaitGroup
	for i, c := range cards {
		wg.Add(1)
		go func(idx int, card Card) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			fullCard, err := GetCardByID(card.ID)
			if err == nil {
				results <- result{idx, fullCard}
			} else {
				results <- result{idx, card}
			}
		}(i, c)
	}

	wg.Wait()
	close(results)

	fullCards := make([]Card, len(cards))
	copy(fullCards, cards)
	for r := range results {
		fullCards[r.index] = r.card
	}
	return fullCards
}

func GetAllCards() ([]Card, error) {
	var cards []Card
	if err := fetchJSON(CardsURL, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}

func GetCardByID(cardID string) (Card, error) {
	var card Card
	if err := fetchJSON(fmt.Sprintf(CardURL, cardID), &card); err != nil {
		return Card{}, err
	}
	card.Image = appendImageSuffix(card.Image, "/high.webp")
	card.SeriesID = card.Set.ID
	card.SeriesName = card.Set.Name
	return card, nil
}

func GetAllSeries() ([]Series, error) {
	var series []Series
	if err := fetchJSON(SetsURL, &series); err != nil {
		return nil, err
	}
	for i := range series {
		series[i].Logo = appendImageSuffix(series[i].Logo, ".webp")
	}
	return series, nil
}

func GetRecentCards(limit int) ([]Card, error) {
	allCards, err := GetAllCards()
	if err != nil {
		return nil, err
	}
	if len(allCards) <= limit {
		return FetchCardDetails(allCards), nil
	}
	return FetchCardDetails(allCards[:limit]), nil
}

func GetSeriesByID(seriesID string) (Series, error) {
	var series Series
	if err := fetchJSON(fmt.Sprintf(SeriesSetURL, seriesID), &series); err != nil {
		return Series{}, err
	}
	series.Logo = appendImageSuffix(series.Logo, ".webp")
	return series, nil
}

func GetCardsBySeries(setID string) ([]Card, error) {
	var setDetail struct {
		Cards []Card `json:"cards"`
	}
	if err := fetchJSON(fmt.Sprintf(SeriesSetURL, setID), &setDetail); err != nil {
		return nil, err
	}
	for i := range setDetail.Cards {
		setDetail.Cards[i].Image = appendImageSuffix(setDetail.Cards[i].Image, "/high.webp")
	}
	return setDetail.Cards, nil
}

func GetCardsByIDs(cardIDs []string) ([]Card, error) {
	if len(cardIDs) == 0 {
		return []Card{}, nil
	}

	minimalCards := make([]Card, len(cardIDs))
	for i, id := range cardIDs {
		minimalCards[i] = Card{ID: id}
	}
	return FetchCardDetails(minimalCards), nil
}

func SearchCards(query SearchQuery) ([]Card, int, error) {
	allCards, err := GetAllCards()
	if err != nil {
		return nil, 0, err
	}

	var filtered []Card
	for _, card := range allCards {
		if query.Query != "" && !strings.Contains(strings.ToLower(card.Name), strings.ToLower(query.Query)) {
			continue
		}
		if query.Category != "" && card.Category != query.Category {
			continue
		}
		if query.Type != "" {
			found := false
			for _, t := range card.Types {
				if t == query.Type {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if query.Rarity != "" && card.Rarity != query.Rarity {
			continue
		}
		filtered = append(filtered, card)
	}

	totalCount := len(filtered)
	start := (query.Page - 1) * query.PageSize
	if start >= totalCount {
		return []Card{}, totalCount, nil
	}
	end := start + query.PageSize
	if end > totalCount {
		end = totalCount
	}

	return FetchCardDetails(filtered[start:end]), totalCount, nil
}

func GetFilterOptions() ([]string, []string, []string, error) {
	cards, err := GetAllCards()
	if err != nil {
		return nil, nil, nil, err
	}

	categories := make(map[string]bool)
	types := make(map[string]bool)
	rarities := make(map[string]bool)

	for _, card := range cards {
		if card.Category != "" {
			categories[card.Category] = true
		}
		for _, t := range card.Types {
			if t != "" {
				types[t] = true
			}
		}
		if card.Rarity != "" {
			rarities[card.Rarity] = true
		}
	}

	toSlice := func(m map[string]bool) []string {
		s := make([]string, 0, len(m))
		for k := range m {
			s = append(s, k)
		}
		return s
	}

	return toSlice(categories), toSlice(types), toSlice(rarities), nil
}

func GetSimilarCards(card Card, limit int) ([]Card, error) {
	allCards, err := GetAllCards()
	if err != nil {
		return nil, err
	}

	var similar []Card
	for _, c := range allCards {
		if c.ID == card.ID || len(similar) >= limit {
			break
		}
		if c.Category != card.Category {
			continue
		}
		if c.Category == "Pokémon" && len(card.Types) > 0 && len(c.Types) > 0 {
			matched := false
			for _, t1 := range card.Types {
				for _, t2 := range c.Types {
					if t1 == t2 {
						matched = true
						break
					}
				}
				if matched {
					break
				}
			}
			if !matched {
				continue
			}
		}
		similar = append(similar, c)
	}

	return FetchCardDetails(similar), nil
}
