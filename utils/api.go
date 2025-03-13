package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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
	body, err := ioutil.ReadAll(resp.Body)
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
