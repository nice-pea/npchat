package sessionn

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func TestNewSession(t *testing.T) {
	t.Run("параметр userID должен быть валидными и не пустыми", func(t *testing.T) {
		session, err := NewSession("invalid", gofakeit.ChromeUserAgent(), StatusNew)
		assert.Zero(t, session)
		assert.ErrorIs(t, err, domain.ErrInvalidID)
	})

	t.Run("параметр name не должен быть пустыми", func(t *testing.T) {
		session, err := NewSession(uuid.NewString(), "", StatusNew)
		assert.Zero(t, session)
		assert.ErrorIs(t, err, ErrSessionNameEmpty)
	})

	t.Run("параметр status должен быть валидными и не пустыми", func(t *testing.T) {
		session, err := NewSession(uuid.NewString(), gofakeit.ChromeUserAgent(), "invalidValue")
		assert.Zero(t, session)
		assert.ErrorIs(t, err, ErrSessionStatusValidate)
	})

	t.Run("status может быть любым", func(t *testing.T) {
		for _, status := range allSessionStatuses {
			session, err := NewSession(uuid.NewString(), gofakeit.ChromeUserAgent(), status)
			assert.NotZero(t, session)
			assert.NoError(t, err)
		}
	})

	t.Run("новой сессии присваивается id, другие свойства равны переданным", func(t *testing.T) {
		userID := uuid.NewString()
		name := "name"
		status := StatusRevoked
		session, err := NewSession(userID, name, status)
		assert.NotZero(t, session)
		assert.NoError(t, err)
		// В id устанавливается случайное значение ID
		assert.NotZero(t, session.ID)
	})

	t.Run("новой сессии создается два токена", func(t *testing.T) {
		session, err := NewSession(uuid.NewString(), gofakeit.ChromeUserAgent(), StatusNew)
		assert.NotZero(t, session)
		assert.NoError(t, err)

		// Токены равны случайными значениями
		assert.NotZero(t, session.AccessToken.Token)
		assert.NotZero(t, session.RefreshToken.Token)

		// Ненулевым сроком жизни
		assert.NotZero(t, session.AccessToken.Expiry)
		assert.NotZero(t, session.RefreshToken.Expiry)
	})
}
