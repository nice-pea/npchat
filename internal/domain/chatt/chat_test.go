package chatt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

func TestNewChat(t *testing.T) {
	t.Run("параметр name должен быть валидными и не пустыми", func(t *testing.T) {
		chat, err := NewChat("\ninvalid\t", uuid.New(), nil)
		assert.Zero(t, chat)
		assert.ErrorIs(t, err, ErrInvalidChatName)
	})

	t.Run("параметр chiefID должен быть валидными и не пустыми", func(t *testing.T) {
		chat, err := NewChat("name", uuid.Nil, nil)
		assert.Zero(t, chat)
		assert.ErrorIs(t, err, ErrInvalidChiefID)
	})

	t.Run("новому чату присваивается id, другие свойства равны переданным", func(t *testing.T) {
		chiefID := uuid.New()
		name := "name"
		chat, err := NewChat(name, chiefID, nil)
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
		chiefID := uuid.New()
		chat, err := NewChat("name", chiefID, nil)
		assert.NotZero(t, chat)
		assert.NoError(t, err)

		// Главный администратор в свойствах чата и участниках
		assert.Len(t, chat.Participants, 1)
		assert.Equal(t, chiefID, chat.Participants[0].UserID)
	})

	t.Run("в новом чате нет приглашений", func(t *testing.T) {
		chat, err := NewChat("name", uuid.New(), nil)
		assert.NotZero(t, chat)
		assert.NoError(t, err)

		// Приглашений нет
		assert.Empty(t, chat.Invitations)
	})

	t.Run("после завершения операции, будут созданы события", func(t *testing.T) {
		// Инициализировать буфер событий
		eventsBuf := new(events.Buffer)

		// Создаем чат
		chat, err := NewChat("name", uuid.New(), eventsBuf)
		assert.NotZero(t, chat)
		assert.NoError(t, err)

		// Проверить список опубликованных событий
		require.Len(t, eventsBuf.Events(), 1)
		// Событие Созданного чата
		chatCreated := eventsBuf.Events()[0].(EventChatCreated)
		// Содержит нужных получателей
		assert.Contains(t, chatCreated.Recipients(), chat.ChiefID)
		// Связано с чатом
		assert.Equal(t, chat.ID, chatCreated.ChatID)
	})
}
