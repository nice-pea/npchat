package repository_tests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func OAuthRepositoryTests(t *testing.T, newRepository func() domain.OAuthRepository) {
	t.Run("SaveToken", func(t *testing.T) {
		t.Run("нельзя сохранить пустой токен", func(t *testing.T) {
			r := newRepository()
			// Сохранить токен
			err := r.SaveToken(domain.OAuthToken{})
			assert.Error(t, err)
		})
		t.Run("сохранить токен со значениями", func(t *testing.T) {
			r := newRepository()
			// Сохранить токен
			err := r.SaveToken(domain.OAuthToken{
				AccessToken:  uuid.NewString(),
				TokenType:    uuid.NewString(),
				RefreshToken: uuid.NewString(),
				Expiry:       time.Now(),
				LinkID:       uuid.NewString(),
				Provider:     uuid.NewString(),
			})
			assert.NoError(t, err)
		})
		t.Run("повторное сохранение будет дублировать запись", func(t *testing.T) {
			r := newRepository()
			// Сохранить токен
			token := rndToken()
			const number = 10
			for range number {
				saveToken(t, r, token)
			}
			// TODO: ListToken
		})
	})
	t.Run("SaveLink", func(t *testing.T) {
		t.Run("нельзя сохранять связь без ID", func(t *testing.T) {
			r := newRepository()
			// Сохранить связь
			savedLink := domain.OAuthLink{
				ID:         "",
				UserID:     "userId",
				ExternalID: "extId",
				Provider:   "provider",
			}
			err := r.SaveLink(savedLink)
			assert.Error(t, err)
		})
		t.Run("сохраненную связь можно прочитать из репозитория", func(t *testing.T) {
			r := newRepository()
			// Сохранить связь
			savedLink := saveLink(t, r, rndLink())
			// Получить связь
			links, err := r.ListLinks(domain.OAuthListLinksFilter{})
			require.NoError(t, err)
			require.Len(t, links, 1)
			assert.Equal(t, savedLink, links[0])
		})
		t.Run("повторное сохранение по одному ID будет обновлять запись", func(t *testing.T) {
			r := newRepository()
			// Много раз сохранить связь
			savedLink := rndLink()
			for range 10 {
				saveLink(t, r, savedLink)
			}
			// Сохранить связь с обновленным полем
			const lastSavedUserID = "someID"
			savedLink.UserID = lastSavedUserID
			saveLink(t, r, savedLink)
			// Получить связь
			links, err := r.ListLinks(domain.OAuthListLinksFilter{})
			require.NoError(t, err)
			require.Len(t, links, 1)
			assert.Equal(t, savedLink, links[0])
		})
	})

	t.Run("ListLinks", func(t *testing.T) {
		t.Run("без фильтра в пустом репозитории", func(t *testing.T) {
			r := newRepository()
			// Список связей
			sessions, err := r.ListLinks(domain.OAuthListLinksFilter{})
			assert.NoError(t, err)
			assert.Empty(t, sessions)
		})
		t.Run("фильтр по ID", func(t *testing.T) {
			r := newRepository()
			// Сохранить связь
			savedLink := saveLink(t, r, rndLink())
			for range 10 {
				saveLink(t, r, rndLink())
			}
			// Список связей
			sessions, err := r.ListLinks(domain.OAuthListLinksFilter{
				ID: savedLink.ID,
			})
			assert.NoError(t, err)
			require.Len(t, sessions, 1)
			assert.Equal(t, savedLink, sessions[0])
		})
		t.Run("фильтр по UserID", func(t *testing.T) {
			r := newRepository()
			// Сохранить связь
			savedLink := saveLink(t, r, rndLink())
			for range 10 {
				saveLink(t, r, rndLink())
			}
			// Список связей
			sessions, err := r.ListLinks(domain.OAuthListLinksFilter{
				UserID: savedLink.UserID,
			})
			assert.NoError(t, err)
			require.Len(t, sessions, 1)
			assert.Equal(t, savedLink, sessions[0])
		})
		t.Run("фильтр по ExternalID", func(t *testing.T) {
			r := newRepository()
			// Сохранить связь
			savedLink := saveLink(t, r, rndLink())
			for range 10 {
				saveLink(t, r, rndLink())
			}
			// Список связей
			sessions, err := r.ListLinks(domain.OAuthListLinksFilter{
				ExternalID: savedLink.ExternalID,
			})
			assert.NoError(t, err)
			require.Len(t, sessions, 1)
			assert.Equal(t, savedLink, sessions[0])
		})
		t.Run("фильтр по Provider", func(t *testing.T) {
			r := newRepository()
			// Сохранить связь
			savedLink := saveLink(t, r, rndLink())
			for range 10 {
				saveLink(t, r, rndLink())
			}
			// Список связей
			sessions, err := r.ListLinks(domain.OAuthListLinksFilter{
				Provider: savedLink.Provider,
			})
			assert.NoError(t, err)
			require.Len(t, sessions, 1)
			assert.Equal(t, savedLink, sessions[0])
		})
		t.Run("фильтр по всем всем полям", func(t *testing.T) {
			r := newRepository()
			// Сохранить связь
			savedLink := saveLink(t, r, rndLink())
			for range 10 {
				saveLink(t, r, rndLink())
			}
			// Список связей
			filter := domain.OAuthListLinksFilter(savedLink)
			sessions, err := r.ListLinks(filter)
			assert.NoError(t, err)
			require.Len(t, sessions, 1)
			assert.Equal(t, savedLink, sessions[0])
		})
	})

}

func saveLink(t *testing.T, r domain.OAuthRepository, l domain.OAuthLink) domain.OAuthLink {
	err := r.SaveLink(l)
	require.NoError(t, err)

	return l
}

func rndLink() domain.OAuthLink {
	return domain.OAuthLink{
		ID:         uuid.NewString(),
		UserID:     uuid.NewString(),
		ExternalID: uuid.NewString(),
		Provider:   uuid.NewString(),
	}
}

func saveToken(t *testing.T, r domain.OAuthRepository, l domain.OAuthToken) domain.OAuthToken {
	err := r.SaveToken(l)
	require.NoError(t, err)

	return l
}

func rndToken() domain.OAuthToken {
	return domain.OAuthToken{
		AccessToken:  uuid.NewString(),
		TokenType:    uuid.NewString(),
		RefreshToken: uuid.NewString(),
		Expiry:       time.Now(),
		LinkID:       uuid.NewString(),
		Provider:     uuid.NewString(),
	}
}
