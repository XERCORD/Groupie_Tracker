<div align="center">

# 🃏 TCGDEX Explorer

**Application web Pokémon TCG — Go · HTML/CSS · API REST**

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![API](https://img.shields.io/badge/API-TCGDEX-EE1515?style=for-the-badge&logo=pokemon&logoColor=white)](https://api.tcgdex.net/)
[![Status](https://img.shields.io/badge/Status-Live-brightgreen?style=for-the-badge)](http://localhost:8000)
[![License](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)](LICENSE)

<br/>

> Explorez, recherchez et sauvegardez vos cartes Pokémon TCG  
> directement depuis l'API officielle TCGDEX — sans framework, 100% Go standard.

<br/>

[🚀 Démarrer](#-installation) · [📖 Documentation](#-routes) · [🔗 API](#-api-tcgdex) · [👤 Auteur](#-auteur)

</div>

---

## 📸 Pages

| Accueil | Recherche | Série |
|:---:|:---:|:---:|
| Cartes récentes + favoris | Filtres nom / type / rareté | Toutes les cartes d'un set |

---

## ✨ Fonctionnalités

```
🔍  Recherche avancée    →  Filtres par nom, type, catégorie, rareté + pagination
📚  Collection de sets   →  Tous les sets Pokémon avec logos, cardCount et tri
🃏  Détails d'une carte  →  Image HD, HP, types, attaques, illustrateur, série
⭐  Favoris persistants  →  Ajout/retrait AJAX, sauvegarde locale JSON
🎴  Cartes similaires    →  Suggestions par type sur la page détail
📄  Pagination           →  Navigation fluide sur toutes les listes
```

---

## ⚡ Stack technique

```go
Language  →  Go 1.23  (zéro dépendance externe)
Serveur   →  net/http  (bibliothèque standard)
Templates →  html/template  (rendu côté serveur)
Frontend  →  HTML5 · CSS3 · JavaScript vanilla
Images    →  TCGDEX CDN  (.webp haute qualité)
Données   →  JSON local  (favoris)
```

---

## 🚀 Installation

### Prérequis

- [Go 1.18+](https://go.dev/dl/)
- [Git](https://git-scm.com/)
- Connexion internet (API TCGDEX en ligne)

### Cloner & lancer

```bash
# Cloner le dépôt
git clone https://github.com/XERCORD/Groupie_Tracker.git
cd Groupie_Tracker

# Lancer le serveur
go run main.go
```

```
✓  Démarrage du serveur TCGDEX Explorer...
✓  API TCGDEX disponible
✓  Serveur démarré sur http://localhost:8000
```

Ouvrir **[http://localhost:8000](http://localhost:8000)** dans le navigateur.

> [!NOTE]
> Aucune dépendance à installer — le projet utilise uniquement la bibliothèque standard Go.

> [!TIP]
> **Port déjà utilisé ?** Trouve et stoppe l'ancienne instance :
> ```powershell
> netstat -ano | findstr ":8000 "
> Stop-Process -Id <PID> -Force
> go run main.go
> ```

---

## 🗂️ Architecture

```
Groupie_Tracker/
│
├── 📄 main.go                    # Serveur HTTP + configuration des routes
├── 📄 go.mod                     # Module Go (aucune dépendance externe)
│
├── 📁 controllers/               # Handlers HTTP
│   ├── home.go                   # Page d'accueil + cartes récentes
│   ├── collection.go             # Liste des sets + page série
│   ├── search.go                 # Recherche avec filtres
│   ├── details.go                # Détails d'une carte
│   ├── favorites.go              # Favoris + API AJAX toggle/clear
│   ├── about.go
│   └── error.go
│
├── 📁 models/                    # Logique métier + appels API
│   ├── api.go                    # Fetch TCGDEX (cards, sets, search…)
│   ├── card.go                   # Structs Card, SearchQuery, CardSet
│   ├── series.go                 # Struct Series, tri, filtres
│   ├── favorites.go              # Persistance JSON des favoris
│   └── helpers.go                # Fonctions utilitaires
│
├── 📁 templates/                 # Templates HTML Go
│   ├── home.html · search.html · collection.html
│   ├── details.html · series.html · favorites.html
│   └── about.html · error.html
│
├── 📁 assets/
│   ├── css/                      # Styles par page (main, home, search…)
│   └── js/favorites.js           # Gestion AJAX des favoris
│
└── 📁 storage/
    └── favorites.json            # Favoris sauvegardés localement
```

---

## 🔌 Routes

### Pages

| Méthode | Route | Description |
|:---:|---|---|
| `GET` | `/accueil` | Page d'accueil — cartes récentes |
| `GET` | `/collection` | Tous les sets — filtrable et triable |
| `GET` | `/series/{id}` | Cartes d'un set avec pagination |
| `GET` | `/carte/{id}` | Détails complets d'une carte |
| `GET` | `/recherche` | Recherche multi-critères |
| `GET` | `/favoris` | Cartes favorites sauvegardées |
| `GET` | `/a-propos` | À propos du projet |

### API REST (AJAX)

| Méthode | Route | Description |
|:---:|---|---|
| `POST` | `/api/favoris/toggle` | Ajouter / retirer une carte des favoris |
| `POST` | `/api/favoris/clear` | Vider tous les favoris |

---

## 🔗 API TCGDEX

Base URL : `https://api.tcgdex.net/v2/fr`

| Endpoint | Utilisation |
|---|---|
| `GET /cards` | Liste minimale de toutes les cartes |
| `GET /cards/{id}` | Détails complets d'une carte (image, HP, attaques…) |
| `GET /sets` | Tous les sets avec cardCount |
| `GET /sets/{id}` | Détails d'un set + cartes avec images incluses |

**Format des images CDN :**
```
https://assets.tcgdex.net/fr/{serie}/{set}/{localId}/high.webp
```

**Optimisation :** les cartes sont récupérées en **parallèle** (goroutines, max 10 simultanées) pour minimiser la latence.

---

## 🛠️ Points techniques notables

<details>
<summary><b>Fetch parallèle des images (goroutines)</b></summary>

Les endpoints liste de l'API ne retournent pas les images. La fonction `FetchCardDetails` récupère les détails de chaque carte en parallèle via des goroutines avec un sémaphore de concurrence :

```go
func FetchCardDetails(cards []Card) []Card {
    sem := make(chan struct{}, 10) // max 10 requêtes simultanées
    results := make(chan result, len(cards))

    for i, card := range cards {
        go func(idx int, c Card) {
            sem <- struct{}{}
            defer func() { <-sem }()
            fullCard, _ := GetCardByID(c.ID)
            results <- result{idx, fullCard}
        }(i, card)
    }
    // ...
}
```
</details>

<details>
<summary><b>Correction du struct Card (champ set imbriqué)</b></summary>

L'API TCGDEX retourne `set` comme un **objet JSON** (pas une string). Le struct a été corrigé pour éviter les erreurs de désérialisation :

```go
type CardSet struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type Card struct {
    Set        CardSet `json:"set"`   // objet imbriqué
    SeriesID   string  `json:"-"`     // peuplé depuis Set.ID
    SeriesName string  `json:"-"`     // peuplé depuis Set.Name
    // ...
}
```
</details>

<details>
<summary><b>Placeholder pour images manquantes (JS)</b></summary>

Certaines cartes anciennes n'ont pas d'image disponible. Un gestionnaire d'erreur JS affiche automatiquement un placeholder 🎴 :

```javascript
document.querySelectorAll('.card-image img').forEach(img => {
    img.addEventListener('error', () => {
        img.style.display = 'none';
        img.closest('.card-image').classList.add('no-image');
    });
});
```
</details>

---

## 📋 Gestion du projet

**Méthodologie :** sprints courts avec priorisation des fonctionnalités

```
Sprint 1  →  Setup Go, routes de base, intégration API
Sprint 2  →  Affichage des cartes + images HD
Sprint 3  →  Recherche avec filtres + pagination
Sprint 4  →  Système de favoris (AJAX + persistance JSON)
Sprint 5  →  Page collection/séries + détails
Sprint 6  →  Fix bugs images, optimisations, README
```

**Ressources utilisées :**
- [Documentation TCGDEX](https://api.tcgdex.net/docs/)
- [Go standard library](https://pkg.go.dev/std)
- [Go net/http](https://pkg.go.dev/net/http)
- [Go html/template](https://pkg.go.dev/html/template)

---

## 👤 Auteur

<div align="center">

**Ji Xerly**  
*Étudiant B1 Informatique — Ynov Campus Aix*

[![GitHub](https://img.shields.io/badge/GitHub-XERCORD-181717?style=for-the-badge&logo=github)](https://github.com/XERCORD)
[![Email](https://img.shields.io/badge/Email-xerly.ji@ynov.com-EA4335?style=for-the-badge&logo=gmail&logoColor=white)](mailto:xerly.ji@ynov.com)

</div>

---

<div align="center">

*Projet Groupie Tracker — Ynov Campus Aix 2025*

</div>
