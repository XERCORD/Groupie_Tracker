function updateFavoriteButton(button, isFavorite) {
  const kind = button.dataset.favKind || 'icon';
  if (kind === 'detail') {
    button.innerHTML = isFavorite ? '❤️ Dans vos favoris' : '🤍 Ajouter aux favoris';
  } else if (kind === 'set') {
    button.innerHTML = isFavorite ? '❤️ Extension en favoris' : '🤍 Ajouter l\'extension aux favoris';
  } else {
    button.innerHTML = isFavorite ? '❤️' : '🤍';
  }
  button.classList.toggle('active', isFavorite);
  button.title = isFavorite ? 'Retirer des favoris' : 'Ajouter aux favoris';
}

function toggleFavoriteRequest(body, button) {
  fetch('/api/favoris/toggle', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
    .then((response) => response.json())
    .then((data) => {
      if (!data.success) {
        showNotification(data.message || 'Une erreur est survenue', 'error');
        return;
      }
      const kind = button.dataset.favKind || 'icon';
      if (kind !== 'set-card') {
        updateFavoriteButton(button, data.isFavorite);
      }

      if (window.location.pathname === '/favoris') {
        if (kind === 'set-card' && !data.isFavorite) {
          const block = button.closest('.favorite-set-card');
          if (block) {
            block.classList.add('removing');
            setTimeout(() => block.remove(), 300);
          }
        }
        if (kind === 'icon' && !data.isFavorite) {
          const card = button.closest('.card');
          if (card) {
            card.classList.add('removing');
            setTimeout(() => {
              card.remove();
              reloadFavoritesIfEmpty();
            }, 300);
          }
        }
        return;
      }

      if (data.isFavorite) {
        showNotification(
          body.setId ? 'Extension ajoutée aux favoris' : 'Carte ajoutée aux favoris',
          'success'
        );
      } else {
        showNotification(
          body.setId ? 'Extension retirée des favoris' : 'Carte retirée des favoris',
          'info'
        );
      }
    })
    .catch((err) => {
      console.error(err);
      showNotification('Une erreur est survenue', 'error');
    });
}

function reloadFavoritesIfEmpty() {
  const cardsLeft = document.querySelectorAll('.cards-grid .card').length;
  const setsLeft = document.querySelectorAll('.favorite-set-card').length;
  if (cardsLeft === 0 && setsLeft === 0) {
    window.location.reload();
  }
}

function clearAllFavorites() {
  if (!confirm('Supprimer toutes les cartes et extensions favorites ?')) return;
  fetch('/api/favoris/clear', { method: 'POST' })
    .then((r) => r.json())
    .then((data) => {
      if (data.success) window.location.reload();
      else showNotification('Une erreur est survenue', 'error');
    })
    .catch(() => showNotification('Une erreur est survenue', 'error'));
}

function showNotification(message, type = 'info') {
  let notification = document.getElementById('notification');
  if (!notification) {
    notification = document.createElement('div');
    notification.id = 'notification';
    document.body.appendChild(notification);
  }
  notification.style.cssText =
    'position:fixed;bottom:20px;right:20px;padding:10px 15px;border-radius:4px;font-size:0.9rem;font-weight:600;z-index:1000;transition:opacity 0.3s';
  const colors = {
    success: ['#4CAF50', '#fff'],
    error: ['#F44336', '#fff'],
    warning: ['#FF9800', '#fff'],
    info: ['#2196F3', '#fff'],
  };
  const [bg, fg] = colors[type] || colors.info;
  notification.style.backgroundColor = bg;
  notification.style.color = fg;
  notification.textContent = message;
  notification.style.opacity = '1';
  setTimeout(() => {
    notification.style.opacity = '0';
    setTimeout(() => notification.remove(), 300);
  }, 3000);
}

function bindFavoriteClick(button) {
  button.addEventListener('click', (e) => {
    e.preventDefault();
    e.stopPropagation();
    const setId = button.getAttribute('data-set-id');
    const cardId = button.getAttribute('data-card-id');
    if (setId) {
      toggleFavoriteRequest({ setId }, button);
    } else if (cardId) {
      toggleFavoriteRequest({ cardId }, button);
    }
  });
}

document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('.favorite-btn').forEach(bindFavoriteClick);
  document.querySelectorAll('.js-favorite-set-toggle').forEach(bindFavoriteClick);

  const clearBtn = document.getElementById('clear-favorites');
  if (clearBtn) clearBtn.addEventListener('click', clearAllFavorites);

  document.querySelectorAll('.card-image img').forEach((img) => {
    img.addEventListener('error', function () {
      this.style.display = 'none';
      const wrap = this.closest('.card-image');
      if (wrap && !wrap.querySelector('.no-image')) {
        wrap.classList.add('no-image');
      }
    });
  });
});

function changePageSize(size) {
  const url = new URL(window.location.href);
  url.searchParams.set('pageSize', size);
  url.searchParams.set('page', '1');
  window.location.href = url.toString();
}
window.changePageSize = changePageSize;
