package chatt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

// TestNewParticipant тестирует создание участника.
func TestNewParticipant(t *testing.T) {
	t.Run("параметр userID должен быть валидным UUID", func(t *testing.T) {
		participant, err := NewParticipant(uuid.Nil)
		assert.Zero(t, participant)
		assert.ErrorIs(t, err, ErrInvalidUserID)
	})

	t.Run("новому участнику присваивается корректный userID", func(t *testing.T) {
		userID := uuid.New()
		participant, err := NewParticipant(userID)
		require.NoError(t, err)
		assert.Equal(t, userID, participant.UserID)
	})
}

// TestChat_AddParticipant тестирует добавление участника в чат.
func TestChat_AddParticipant(t *testing.T) {
	t.Run("добавление участника в чат", func(t *testing.T) {
		// Создаем чат
		chat, err := NewChat("test chat", uuid.New(), nil)
		require.NoError(t, err)

		// Создаем нового участника
		userID := uuid.New()
		participant, err := NewParticipant(userID)
		require.NoError(t, err)

		// Добавляем участника
		err = chat.AddParticipant(participant, nil)
		require.NoError(t, err)

		// Проверяем наличие участника
		assert.True(t, chat.HasParticipant(userID))
	})

	t.Run("нельзя добавить уже существующего участника", func(t *testing.T) {
		// Создаем чат
		chat, err := NewChat("test chat", uuid.New(), nil)
		require.NoError(t, err)

		// Создаем участника
		userID := uuid.New()
		participant, err := NewParticipant(userID)
		require.NoError(t, err)

		// Добавляем первый раз
		err = chat.AddParticipant(participant, nil)
		require.NoError(t, err)

		// Пробуем добавить повторно
		err = chat.AddParticipant(participant, nil)
		assert.ErrorIs(t, err, ErrParticipantExists)
	})

	t.Run("нельзя добавить если есть приглашение к этому участнику", func(t *testing.T) {
		// Создаем чат
		chat, err := NewChat("test chat", uuid.New(), nil)
		require.NoError(t, err)

		userID := uuid.New()

		// Создаем и добавляем приглашение
		inv, err := NewInvitation(chat.ChiefID, userID)
		require.NoError(t, err)
		err = chat.AddInvitation(inv, nil)
		require.NoError(t, err)

		// Создаем и добавляем участника
		participant, err := NewParticipant(userID)
		require.NoError(t, err)
		err = chat.AddParticipant(participant, nil)
		assert.ErrorIs(t, err, ErrUserIsAlreadyInvited)
	})

	t.Run("после завершения операции, будут созданы события", func(t *testing.T) {
		// Создаем чат
		chiefID := uuid.New()
		chat, err := NewChat("test chat", chiefID, nil)
		require.NoError(t, err)

		// Инициализировать буфер событий
		eventsBuf := new(events.Buffer)

		// Создаем и добавляем участника
		userID := uuid.New()
		participant, err := NewParticipant(userID)
		require.NoError(t, err)
		err = chat.AddParticipant(participant, eventsBuf)
		require.NoError(t, err)

		// Проверить список опубликованных событий
		require.Len(t, eventsBuf.Events(), 1)
		// Событие Удаленного
		participantAdded := eventsBuf.Events()[0].(EventParticipantAdded)
		// Содержит нужных получателей
		assert.Contains(t, participantAdded.Recipients(), chat.ChiefID)
		assert.Contains(t, participantAdded.Recipients(), participant.UserID)
		// Связано с чатом
		assert.Equal(t, chat.ID, participantAdded.ChatID)
		// Содержит нужного участника
		assert.Equal(t, participant, participantAdded.Participant)
	})
}

// TestChat_RemoveParticipant тестирует удаление участника из чата.
func TestChat_RemoveParticipant(t *testing.T) {
	t.Run("удаление участника из чата", func(t *testing.T) {
		// Создаем чат
		chiefID := uuid.New()
		chat, err := NewChat("test chat", chiefID, nil)
		require.NoError(t, err)

		// Создаем и добавляем участника
		userID := uuid.New()
		participant, err := NewParticipant(userID)
		require.NoError(t, err)
		err = chat.AddParticipant(participant, nil)
		require.NoError(t, err)

		// Удаляем участника
		err = chat.RemoveParticipant(userID, nil)
		require.NoError(t, err)

		// Проверяем, что участник удален
		assert.False(t, chat.HasParticipant(userID))
	})

	t.Run("нельзя удалить несуществующего участника", func(t *testing.T) {
		// Создаем чат
		chat, err := NewChat("test chat", uuid.New(), nil)
		require.NoError(t, err)

		// Удаляем участника (несуществующего)
		err = chat.RemoveParticipant(uuid.New(), nil)
		assert.ErrorIs(t, err, ErrParticipantNotExists)
	})

	t.Run("нельзя удалить главного администратора", func(t *testing.T) {
		// Создаем чат
		chiefID := uuid.New()
		chat, err := NewChat("test chat", chiefID, nil)
		require.NoError(t, err)

		// Удаляем участника (главного администратора)
		err = chat.RemoveParticipant(chiefID, nil)
		assert.ErrorIs(t, err, ErrCannotRemoveChief)
	})

	t.Run("после завершения операции, будут созданы события", func(t *testing.T) {
		// Создаем чат
		chiefID := uuid.New()
		chat, err := NewChat("test chat", chiefID, nil)
		require.NoError(t, err)

		// Создаем и добавляем участника
		userID := uuid.New()
		participant, err := NewParticipant(userID)
		require.NoError(t, err)
		err = chat.AddParticipant(participant, nil)
		require.NoError(t, err)

		// Инициализировать буфер событий
		eventsBuf := new(events.Buffer)

		// Удаляем участника
		err = chat.RemoveParticipant(userID, eventsBuf)
		require.NoError(t, err)

		// Проверить список опубликованных событий
		require.Len(t, eventsBuf.Events(), 1)
		// Событие Удаленного
		participantRemoved := eventsBuf.Events()[0].(EventParticipantRemoved)
		// Содержит нужных получателей
		assert.Contains(t, participantRemoved.Recipients(), chat.ChiefID)
		assert.Contains(t, participantRemoved.Recipients(), participant.UserID)
		// Связано с чатом
		assert.Equal(t, chat.ID, participantRemoved.ChatID)
		// Содержит нужного участника
		assert.Equal(t, participant, participantRemoved.Participant)
	})
}
