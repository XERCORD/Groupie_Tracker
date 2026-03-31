package models

import (
	"encoding/json"
	"os"
	"sync"
)

type Favorites struct {
	CardIDs []string `json:"favorites"`
	SetIDs  []string `json:"favoriteSets,omitempty"`
	mu      sync.Mutex
}

const favoritesFilePath = "./storage/favorites.json"

var (
	favoritesInstance *Favorites
	once              sync.Once
)

func GetFavorites() *Favorites {
	once.Do(func() {
		favoritesInstance = &Favorites{CardIDs: []string{}, SetIDs: []string{}}
		favoritesInstance.Load()
	})
	return favoritesInstance
}

func (f *Favorites) Load() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, err := os.Stat(favoritesFilePath); os.IsNotExist(err) {
		os.MkdirAll("./storage", 0755)
		f.CardIDs = []string{}
		f.SetIDs = []string{}
		return f.saveWithoutLock()
	}

	data, err := os.ReadFile(favoritesFilePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		f.CardIDs = []string{}
		f.SetIDs = []string{}
		return nil
	}

	if err := json.Unmarshal(data, f); err != nil {
		return err
	}
	if f.CardIDs == nil {
		f.CardIDs = []string{}
	}
	if f.SetIDs == nil {
		f.SetIDs = []string{}
	}
	return nil
}

func (f *Favorites) Save() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.saveWithoutLock()
}

func (f *Favorites) saveWithoutLock() error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return os.WriteFile(favoritesFilePath, data, 0644)
}

func (f *Favorites) Add(cardID string) error {
	f.mu.Lock()
	for _, id := range f.CardIDs {
		if id == cardID {
			f.mu.Unlock()
			return nil
		}
	}
	f.CardIDs = append(f.CardIDs, cardID)
	f.mu.Unlock()
	return f.Save()
}

func (f *Favorites) Remove(cardID string) error {
	f.mu.Lock()
	index := -1
	for i, id := range f.CardIDs {
		if id == cardID {
			index = i
			break
		}
	}
	if index >= 0 {
		f.CardIDs = append(f.CardIDs[:index], f.CardIDs[index+1:]...)
	}
	f.mu.Unlock()
	return f.Save()
}

func (f *Favorites) Toggle(cardID string) (bool, error) {
	isFavorite := f.Contains(cardID)
	var err error
	if isFavorite {
		err = f.Remove(cardID)
	} else {
		err = f.Add(cardID)
	}
	return !isFavorite, err
}

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

func (f *Favorites) GetAll() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	result := make([]string, len(f.CardIDs))
	copy(result, f.CardIDs)
	return result
}

func (f *Favorites) AddSet(setID string) error {
	f.mu.Lock()
	for _, id := range f.SetIDs {
		if id == setID {
			f.mu.Unlock()
			return nil
		}
	}
	f.SetIDs = append(f.SetIDs, setID)
	f.mu.Unlock()
	return f.Save()
}

func (f *Favorites) RemoveSet(setID string) error {
	f.mu.Lock()
	index := -1
	for i, id := range f.SetIDs {
		if id == setID {
			index = i
			break
		}
	}
	if index >= 0 {
		f.SetIDs = append(f.SetIDs[:index], f.SetIDs[index+1:]...)
	}
	f.mu.Unlock()
	return f.Save()
}

func (f *Favorites) ToggleSet(setID string) (bool, error) {
	isFav := f.ContainsSet(setID)
	var err error
	if isFav {
		err = f.RemoveSet(setID)
	} else {
		err = f.AddSet(setID)
	}
	return !isFav, err
}

func (f *Favorites) ContainsSet(setID string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, id := range f.SetIDs {
		if id == setID {
			return true
		}
	}
	return false
}

func (f *Favorites) GetAllSets() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]string, len(f.SetIDs))
	copy(out, f.SetIDs)
	return out
}

func (f *Favorites) Clear() error {
	f.mu.Lock()
	f.CardIDs = []string{}
	f.SetIDs = []string{}
	f.mu.Unlock()
	return f.Save()
}

func (f *Favorites) Count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.CardIDs)
}
