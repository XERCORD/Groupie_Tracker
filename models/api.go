package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
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

type cacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

var (
	cacheMu         sync.Mutex
	allCardsCache   *cacheEntry
	serieBlocksCache *cacheEntry
	cardCache       = make(map[string]*cacheEntry)
	cacheTTL        = 5 * time.Minute
)

type tcgSeriesListResponse struct {
	Value []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"value"`
}

type serieBlockData struct {
	ByID    map[string]string
	IDsDesc []string
}

func getSerieBlockData() (*serieBlockData, error) {
	cacheMu.Lock()
	if serieBlocksCache != nil && time.Now().Before(serieBlocksCache.expiresAt) {
		d := serieBlocksCache.data.(*serieBlockData)
		cacheMu.Unlock()
		return d, nil
	}
	cacheMu.Unlock()

	var raw tcgSeriesListResponse
	if err := fetchJSON(SeriesURL, &raw); err != nil {
		return nil, err
	}

	byID := make(map[string]string)
	ids := make([]string, 0, len(raw.Value))
	for _, v := range raw.Value {
		if v.ID == "" {
			continue
		}
		byID[v.ID] = v.Name
		ids = append(ids, v.ID)
	}
	sort.Slice(ids, func(i, j int) bool {
		if len(ids[i]) != len(ids[j]) {
			return len(ids[i]) > len(ids[j])
		}
		return ids[i] < ids[j]
	})

	d := &serieBlockData{ByID: byID, IDsDesc: ids}
	cacheMu.Lock()
	serieBlocksCache = &cacheEntry{data: d, expiresAt: time.Now().Add(cacheTTL)}
	cacheMu.Unlock()
	return d, nil
}

func serieIDForSet(setID string, data *serieBlockData) string {
	if data == nil {
		return ""
	}
	for _, id := range data.IDsDesc {
		if strings.HasPrefix(setID, id) {
			return id
		}
	}
	return ""
}

func GetSerieOptions() ([]SerieOption, error) {
	data, err := getSerieBlockData()
	if err != nil {
		return nil, err
	}
	opts := make([]SerieOption, 0, len(data.ByID))
	for id, name := range data.ByID {
		opts = append(opts, SerieOption{ID: id, Name: name})
	}
	sort.Slice(opts, func(i, j int) bool {
		return strings.ToLower(opts[i].Name) < strings.ToLower(opts[j].Name)
	})
	return opts, nil
}

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
	cacheMu.Lock()
	if allCardsCache != nil && time.Now().Before(allCardsCache.expiresAt) {
		cards := allCardsCache.data.([]Card)
		cacheMu.Unlock()
		return cards, nil
	}
	cacheMu.Unlock()

	var cards []Card
	if err := fetchJSON(CardsURL, &cards); err != nil {
		return nil, err
	}

	cacheMu.Lock()
	allCardsCache = &cacheEntry{data: cards, expiresAt: time.Now().Add(cacheTTL)}
	cacheMu.Unlock()

	return cards, nil
}

func GetCardByID(cardID string) (Card, error) {
	cacheMu.Lock()
	if entry, ok := cardCache[cardID]; ok && time.Now().Before(entry.expiresAt) {
		card := entry.data.(Card)
		cacheMu.Unlock()
		return card, nil
	}
	cacheMu.Unlock()

	var card Card
	if err := fetchJSON(fmt.Sprintf(CardURL, cardID), &card); err != nil {
		return Card{}, err
	}
	card.Image = appendImageSuffix(card.Image, "/high.webp")
	card.SeriesID = card.Set.ID
	card.SeriesName = card.Set.Name

	cacheMu.Lock()
	cardCache[cardID] = &cacheEntry{data: card, expiresAt: time.Now().Add(cacheTTL)}
	cacheMu.Unlock()

	return card, nil
}

func GetAllSeries() ([]Series, error) {
	var series []Series
	if err := fetchJSON(SetsURL, &series); err != nil {
		return nil, err
	}
	blocks, err := getSerieBlockData()
	if err != nil {
		blocks = &serieBlockData{ByID: map[string]string{}, IDsDesc: nil}
	}
	for i := range series {
		series[i].Logo = appendImageSuffix(series[i].Logo, ".webp")
		sid := serieIDForSet(series[i].ID, blocks)
		series[i].SerieID = sid
		if sid != "" {
			series[i].SerieName = blocks.ByID[sid]
		}
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
	blocks, _ := getSerieBlockData()
	if blocks != nil {
		sid := serieIDForSet(series.ID, blocks)
		series.SerieID = sid
		if sid != "" {
			series.SerieName = blocks.ByID[sid]
		}
	}
	return series, nil
}

func GetSetsByIDs(setIDs []string) ([]Series, error) {
	if len(setIDs) == 0 {
		return nil, nil
	}
	all, err := GetAllSeries()
	if err != nil {
		return nil, err
	}
	byID := make(map[string]Series, len(all))
	for _, s := range all {
		byID[s.ID] = s
	}
	out := make([]Series, 0, len(setIDs))
	for _, id := range setIDs {
		if s, ok := byID[id]; ok {
			out = append(out, s)
		}
	}
	return out, nil
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
