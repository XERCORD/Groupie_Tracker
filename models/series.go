package models

import (
	"sort"
	"strings"
)

type Series struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ReleaseDate string    `json:"releaseDate"`
	CardCount   CardCount `json:"cardCount"`
	Logo        string    `json:"logo,omitempty"`
}

type CardCount struct {
	Total    int `json:"total"`
	Official int `json:"official"`
}

func SortSeries(series []Series, by string) []Series {
	result := make([]Series, len(series))
	copy(result, series)

	switch by {
	case "date-desc":
		sort.Slice(result, func(i, j int) bool { return result[i].ReleaseDate > result[j].ReleaseDate })
	case "date-asc":
		sort.Slice(result, func(i, j int) bool { return result[i].ReleaseDate < result[j].ReleaseDate })
	case "cards-desc":
		sort.Slice(result, func(i, j int) bool { return result[i].CardCount.Total > result[j].CardCount.Total })
	case "cards-asc":
		sort.Slice(result, func(i, j int) bool { return result[i].CardCount.Total < result[j].CardCount.Total })
	default:
		sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	}

	return result
}

func FilterSeries(series []Series, search string) []Series {
	if search == "" {
		return series
	}
	lower := strings.ToLower(search)
	var result []Series
	for _, s := range series {
		if strings.Contains(strings.ToLower(s.Name), lower) {
			result = append(result, s)
		}
	}
	return result
}
