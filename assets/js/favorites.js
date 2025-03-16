/**
 * Script pour gérer les favoris
 */

// Fonction pour ajouter/retirer une carte des favoris
function toggleFavorite(cardId, button) {
    // Appel AJAX pour modifier les favoris
    fetch('/api/favoris/toggle', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ cardId: cardId }),
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        // Mettre à jour l'interface utilisateur
        if (data.isFavorite) {
          // Ajouté aux favoris
          button.innerHTML = '❤️';
          button.classList.add('active');
          button.title = 'Retirer des favoris';
          
          // Afficher une notification
          showNotification('Carte ajoutée aux favoris', 'success');
        } else {
          // Retiré des favoris
          button.innerHTML = '🤍';
          button.classList.remove('active');
          button.title = 'Ajouter aux favoris';
          
          // Si on est sur la page des favoris, supprimer la carte visuellement
          if (window.location.pathname === '/favoris') {
            const card = button.closest('.card');
            if (card) {
              card.classList.add('removing');
              setTimeout(() => {
                card.remove();
                
                // Vérifier s'il reste des cartes
                const remainingCards = document.querySelectorAll('.cards-grid .card');
                if (remainingCards.length === 0) {
                  // Recharger la page pour afficher le message "Aucun favori"
                  window.location.reload();
                }
              }, 300);
            }
          }
          
          // Afficher une notification
          showNotification('Carte retirée des favoris', 'info');
        }
      } else {
        // Erreur
        showNotification('Une erreur est survenue', 'error');
      }
    })
    .catch(error => {
      console.error('Erreur:', error);
      showNotification('Une erreur est survenue', 'error');
    });
  }
  
  // Fonction pour changer la taille de la page dans les résultats de recherche
  function changePageSize(size) {
    // Récupérer l'URL actuelle
    const url = new URL(window.location.href);
    
    // Mettre à jour le paramètre pageSize
    url.searchParams.set('pageSize', size);
    
    // Revenir à la page 1
    url.searchParams.set('page', '1');
    
    // Rediriger vers la nouvelle URL
    window.location.href = url.toString();
  }
  
  // Fonction pour supprimer tous les favoris
  function clearAllFavorites() {
    if (confirm('Êtes-vous sûr de vouloir supprimer tous vos favoris ?')) {
      fetch('/api/favoris/clear', {
        method: 'POST',
      })
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          // Recharger la page
          window.location.reload();
        } else {
          showNotification('Une erreur est survenue', 'error');
        }
      })
      .catch(error => {
        console.error('Erreur:', error);
        showNotification('Une erreur est survenue', 'error');
      });
    }
  }
  
  // Fonction pour afficher des notifications
  function showNotification(message, type = 'info') {
    // Créer l'élément de notification s'il n'existe pas déjà
    let notification = document.getElementById('notification');
    if (!notification) {
      notification = document.createElement('div');
      notification.id = 'notification';
      document.body.appendChild(notification);
    }
    
    // Définir le style de base
    notification.style.position = 'fixed';
    notification.style.bottom = '20px';
    notification.style.right = '20px';
    notification.style.padding = '10px 15px';
    notification.style.borderRadius = '4px';
    notification.style.fontSize = '0.9rem';
    notification.style.fontWeight = '600';
    notification.style.zIndex = '1000';
    notification.style.transition = 'opacity 0.3s';
    
    // Définir la couleur selon le type
    switch (type) {
      case 'success':
        notification.style.backgroundColor = '#4CAF50';
        notification.style.color = '#fff';
        break;
      case 'error':
        notification.style.backgroundColor = '#F44336';
        notification.style.color = '#fff';
        break;
      case 'warning':
        notification.style.backgroundColor = '#FF9800';
        notification.style.color = '#fff';
        break;
      default:
        notification.style.backgroundColor = '#2196F3';
        notification.style.color = '#fff';
    }
    
    // Définir le message
    notification.textContent = message;
    
    // Afficher la notification
    notification.style.opacity = '1';
    
    // Cacher la notification après 3 secondes
    setTimeout(() => {
      notification.style.opacity = '0';
      setTimeout(() => {
        notification.remove();
      }, 300);
    }, 3000);
  }
  
  // Vérifier si la page est chargée
  document.addEventListener('DOMContentLoaded', () => {
    // Initialiser les boutons de favoris
    const favoriteButtons = document.querySelectorAll('.favorite-btn');
    favoriteButtons.forEach(button => {
      const cardId = button.getAttribute('data-card-id');
      button.addEventListener('click', () => toggleFavorite(cardId, button));
    });
    
    // Initialiser le bouton pour supprimer tous les favoris (si présent)
    const clearButton = document.getElementById('clear-favorites');
    if (clearButton) {
      clearButton.addEventListener('click', clearAllFavorites);
    }
  });