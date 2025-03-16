package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// Favorites représente la liste des cartes favorites
type Favorites struct {
	CardIDs []string `json:"favorites"`
	mu      sync.Mutex
}

// Chemin du fichier pour la persistance des favoris
const favoritesFilePath = "./storage/favorites.json"

// Instance globale des favoris (singleton)
var (
	favoritesInstance *Favorites
	once              sync.Once
)

// GetFavorites retourne l'instance unique des favoris
func GetFavorites() *Favorites {
	once.Do(func() {
		favoritesInstance = &Favorites{
			CardIDs: []string{},
		}

		// Charge les favoris depuis le fichier
		favoritesInstance.Load()
	})

	return favoritesInstance
}

// Load charge les favoris depuis le fichier
func (f *Favorites) Load() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Vérifie si le fichier existe
	if _, err := os.Stat(favoritesFilePath); os.IsNotExist(err) {
		// Crée le dossier storage s'il n'existe pas
		os.MkdirAll("./storage", 0755)

		// Initialise avec un tableau vide
		f.CardIDs = []string{}

		// Sauvegarde sans verrouillage (nous avons déjà le verrou)
		return f.saveWithoutLock()
	}

	// Lit le fichier
	data, err := ioutil.ReadFile(favoritesFilePath)
	if err != nil {
		return err
	}

	// Si le fichier est vide, initialise avec un tableau vide
	if len(data) == 0 {
		f.CardIDs = []string{}
		return nil
	}

	// Désérialise le JSON
	return json.Unmarshal(data, f)
}

// Save sauvegarde les favoris dans le fichier
func (f *Favorites) Save() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.saveWithoutLock()
}

// saveWithoutLock sauvegarde les favoris sans verrouiller le mutex
// à utiliser uniquement lorsque le mutex est déjà verrouillé
func (f *Favorites) saveWithoutLock() error {
	// Sérialise en JSON
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}

	// Écrit dans le fichier
	return ioutil.WriteFile(favoritesFilePath, data, 0644)
}

// Add ajoute une carte aux favoris
func (f *Favorites) Add(cardID string) error {
	f.mu.Lock()

	// Vérifie si la carte est déjà dans les favoris
	for _, id := range f.CardIDs {
		if id == cardID {
			f.mu.Unlock()
			return nil // Déjà présente, rien à faire
		}
	}

	// Ajoute l'ID de la carte
	f.CardIDs = append(f.CardIDs, cardID)
	f.mu.Unlock()

	// Sauvegarde les changements
	return f.Save()
}

// Remove retire une carte des favoris
func (f *Favorites) Remove(cardID string) error {
	f.mu.Lock()

	// Recherche l'index de la carte
	index := -1
	for i, id := range f.CardIDs {
		if id == cardID {
			index = i
			break
		}
	}

	// Si trouvé, retire la carte
	if index >= 0 {
		f.CardIDs = append(f.CardIDs[:index], f.CardIDs[index+1:]...)
	}

	f.mu.Unlock()

	// Sauvegarde les changements
	return f.Save()
}

// Toggle ajoute ou retire une carte des favoris
func (f *Favorites) Toggle(cardID string) (bool, error) {
	// Vérifie si la carte est déjà dans les favoris
	isFavorite := f.Contains(cardID)

	var err error
	if isFavorite {
		// Si oui, la retire
		err = f.Remove(cardID)
	} else {
		// Sinon, l'ajoute
		err = f.Add(cardID)
	}

	// Retourne le nouvel état (inversé) et l'erreur éventuelle
	return !isFavorite, err
}

// Contains vérifie si une carte est dans les favoris
func (f *Favorites) Contains(cardID string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, id := range f.CardIDs {
		if id == cardID {
			return true
		}
	}

	return false
}

// GetAll retourne une copie de tous les IDs de cartes favoris
func (f *Favorites) GetAll() []string {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Crée une copie pour éviter les problèmes de concurrence
	result := make([]string, len(f.CardIDs))
	copy(result, f.CardIDs)

	return result
}

// Clear supprime toutes les cartes des favoris
func (f *Favorites) Clear() error {
	f.mu.Lock()
	f.CardIDs = []string{}
	f.mu.Unlock()

	// Sauvegarde les changements
	return f.Save()
}

// Count retourne le nombre de cartes favorites
func (f *Favorites) Count() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return len(f.CardIDs)
}
