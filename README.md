# TCGDEX Explorer

TCGDEX Explorer est une application permettant d'explorer les cartes du jeu de cartes à collectionner via l'API publique TCGDEX.  
L'objectif est de fournir une interface intuitive permettant de rechercher des cartes, consulter des collections et gérer une liste de favoris.

## 📌 Fonctionnalités

- 🔍 Recherche de cartes avec filtres  
- 📚 Affichage des collections et sets disponibles  
- 🃏 Consultation des détails d'une carte spécifique  
- ⭐ Gestion des favoris (ajout/suppression)  

---

## 🚀 Installation et exécution

### 1️⃣ Prérequis
- **Go** (1.18 ou supérieur)
- **Git** (pour cloner le projet)

### 2️⃣ Cloner le projet
1. Cloner le dépôt :
   ```bash
   git clone https://github.com/XERCORD/Groupie_Tracker.git
2. Lancer le serveur :
   ```bash
   go run main.go
3. Accéder à l'application
  Ouvrez le navigateur Chrome et entrez l'adresse : http://localhost:8000

## 📌 Routes implémentées

| Route               | Méthode | Description                                      |
|---------------------|---------|--------------------------------------------------|
| `/accueil`         | GET     | Page d'accueil                                  |
| `/recherche`       | GET     | Page de recherche de cartes                     |
| `/collection`      | GET     | Affiche tous les sets disponibles               |
| `/card`            | GET     | Affiche les détails d'une carte spécifique      |
| `/favorites`       | GET     | Affiche la liste des favoris                    |
| `/favorite/toggle` | POST    | Ajoute ou retire une carte des favoris          |


## 🔗 API utilisée

L'application utilise l'API publique [TCGDEX](https://api.tcgdex.net/). Voici les endpoints exploités :

- **Obtenir toutes les cartes** : `GET https://api.tcgdex.net/v2/fr/cards`
- **Obtenir un set** : `GET https://api.tcgdex.net/v2/fr/{id}`
- **Obtenir toutes les séries** : `GET https://api.tcgdex.net/v2/fr/series`

## 📖 À propos du projet

### ❓ FAQ – Gestion du projet

#### 🔹 Comment avez-vous décomposé le projet ? Quelles ont été les phases clé ?
Le projet a été divisé en plusieurs phases :  

1️⃣ **Analyse des besoins** : Identification des fonctionnalités essentielles et définition des contraintes techniques.  
2️⃣ **Exploration de l'API** : Étude de l'API TCGDEX et tests des endpoints pour comprendre les données disponibles.  
3️⃣ **Conception de l'interface** : Création de wireframes avec Canva pour structurer l'affichage des cartes et collections.  
4️⃣ **Développement backend** : Implémentation des routes et de la logique métier en Go pour interagir avec l'API.  
5️⃣ **Développement frontend** : Intégration des templates HTML et du CSS pour afficher les résultats dynamiquement.  
6️⃣ **Tests et débogage** : Vérification du bon fonctionnement, correction des bugs et optimisation du code.  

#### 🔹 Comment avez-vous réparti les tâches ?
Le projet étant **individuel**, j'ai organisé mon travail de manière agile en utilisant un **Trello** pour suivre mes tâches et priorités.  
J'ai divisé mon travail en **sprints courts**, chaque sprint étant consacré à une fonctionnalité spécifique (recherche, favoris, affichage des détails…).

#### 🔹 Comment avez-vous géré votre temps ?
J'ai utilisé la **méthode Pomodoro** pour structurer mon temps de travail en sessions de 25 minutes avec des pauses courtes.  
J'ai également défini des **priorités** selon l'importance des fonctionnalités :

1️⃣ **Fonctionnalités essentielles** :  
   - Affichage des cartes  
   - Recherche avec filtres  
   - Système de favoris  

2️⃣ **Fonctionnalités secondaires** :  
   - Pagination des résultats  
   - Optimisation de l'affichage  

#### 🔹 Quelle stratégie avez-vous adoptée pour vous documenter ?
Pour comprendre et résoudre les problèmes techniques, j’ai utilisé plusieurs sources de documentation :  

- 📜 **Documentation officielle** de [l'API TCGDEX](https://api.tcgdex.net/docs/).  
- 📘 **Documentation Go** pour la gestion des routes et des templates HTML.  
- 🎥 **Tutoriels en ligne** sur le développement web avec Go et la consommation d'API REST.  
- 💬 **Forums et communautés**  pour obtenir des solutions à des problèmes spécifiques.  

---

## 👤 Auteur

**Ji Xerly**  
Développeur du projet **TCGDEX Explorer**.  
📧 Contact : [xerly.ji@ynov.com](mailto:xerly.ji@ynov.com)  
🔗 GitHub : [github.com/XERCORD](https://github.com/XERCORD)

