package chatt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewChat(t *testing.T) {
	t.Run("параметр name не должен быть пустым", func(t *testing.T) {
		chat, err := NewChat("", uuid.NewString())
		assert.Zero(t, chat)
		assert.ErrorIs(t, err, ErrInvalidChatName)
	})

	t.Run("параметр name должен быть валидными и не пустыми", func(t *testing.T) {
		chat, err := NewChat("\ninvalid\t", uuid.NewString())
		assert.Zero(t, chat)
		assert.ErrorIs(t, err, ErrInvalidChatName)
	})

	t.Run("параметр chiefID не должен быть пустым", func(t *testing.T) {
		chat, err := NewChat("name", "")
		assert.Zero(t, chat)
		assert.ErrorIs(t, err, ErrInvalidChiefID)
	})

	t.Run("параметр chiefID должен быть валидным", func(t *testing.T) {
		chat, err := NewChat("name", "invalid")
		assert.Zero(t, chat)
		assert.ErrorIs(t, err, ErrInvalidChiefID)
	})

	t.Run("новому чату присваивается id, другие свойства равны переданным", func(t *testing.T) {
		chiefID := uuid.NewString()
		name := "name"
		chat, err := NewChat(name, chiefID)
		assert.NotZero(t, chat)
		assert.NoError(t, err)

		// В id устанавливается случайное значение ID
		assert.NotZero(t, chat.ID)
		// Главный администратор из параметров
		assert.Equal(t, chiefID, chat.ChiefID)
		// Название чата из параметров
		assert.Equal(t, name, chat.Name)
	})

	t.Run("в новом чате создается главный администратор", func(t *testing.T) {
		chiefID := uuid.NewString()
		chat, err := NewChat("name", chiefID)
		assert.NotZero(t, chat)
		assert.NoError(t, err)

		// Главный администратор в свойствах чата и участниках
		assert.Len(t, chat.Participants, 1)
		assert.Equal(t, chiefID, chat.Participants[0].UserID)
	})

	t.Run("в новом чате нет приглашений", func(t *testing.T) {
		chat, err := NewChat("name", uuid.NewString())
		assert.NotZero(t, chat)
		assert.NoError(t, err)

		// Приглашений нет
		assert.Empty(t, chat.Invitations)
	})

}
