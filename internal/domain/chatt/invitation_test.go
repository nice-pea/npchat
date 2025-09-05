package chatt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

// TestNewInvitation тестирует создание приглашения.
func TestNewInvitation(t *testing.T) {
	t.Run("параметр subjectID не должен быть пустым", func(t *testing.T) {
		inv, err := NewInvitation(uuid.Nil, uuid.New())
		assert.Zero(t, inv)
		assert.Error(t, err)
	})

	t.Run("параметр recipientID не должен быть пустым", func(t *testing.T) {
		inv, err := NewInvitation(uuid.New(), uuid.Nil)
		assert.Zero(t, inv)
		assert.Error(t, err)
	})

	t.Run("subjectID и recipientID не могут быть одинаковыми", func(t *testing.T) {
		id := uuid.New()
		inv, err := NewInvitation(id, id)
		assert.Zero(t, inv)
		assert.ErrorIs(t, err, ErrSubjectAndRecipientMustBeDifferent)
	})

	t.Run("создание валидного приглашения", func(t *testing.T) {
		subjectID := uuid.New()
		recipientID := uuid.New()
		inv, err := NewInvitation(subjectID, recipientID)
		assert.NotZero(t, inv)
		assert.NoError(t, err)
		// В id устанавливается случайное значение ID
		assert.NotEmpty(t, inv.ID)
		// Свойства из параметров конструктора
		assert.Equal(t, subjectID, inv.SubjectID)
		assert.Equal(t, recipientID, inv.RecipientID)
	})
}

// TestNewInvitation тестирует добавление приглашения в чат.
func TestChat_AddInvitation(t *testing.T) {
	t.Run("нельзя добавить приглашение от не участника чата", func(t *testing.T) {
		// Создать чат
		chat, err := NewChat("chatName", uuid.New(), nil)
		require.NoError(t, err)

		// Создать и добавить первое приглашение
		inv, err := NewInvitation(uuid.New(), uuid.New())
		require.NoError(t, err)
		err = chat.AddInvitation(inv, nil)
		assert.ErrorIs(t, err, ErrSubjectIsNotMember)
	})

	t.Run("нельзя пригласить существующего участника", func(t *testing.T) {
		// Создать чат
		chief := uuid.New()
		chat, err := NewChat("test", chief, nil)
		require.NoError(t, err)

		// Создать и добавить участника
		p, err := NewParticipant(uuid.New())
		require.NoError(t, err)
		err = chat.AddParticipant(p, nil)
		require.NoError(t, err)

		// Добавить приглашение пользователю, который уже участник
		inv, err := NewInvitation(chief, p.UserID)
		require.NoError(t, err)
		err = chat.AddInvitation(inv, nil)
		assert.ErrorIs(t, err, ErrParticipantExists)
	})

	t.Run("нельзя пригласить уже приглашенного пользователя", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.New()
		recipientID := uuid.New()
		chat, err := NewChat("chatName", chiefID, nil)
		require.NoError(t, err)

		// Создать и добавить первое приглашение
		inv1, err := NewInvitation(chiefID, recipientID)
		require.NoError(t, err)
		err = chat.AddInvitation(inv1, nil)
		require.NoError(t, err)

		//  Создать и добавить второе приглашение
		inv2, err := NewInvitation(chiefID, recipientID)
		require.NoError(t, err)
		err = chat.AddInvitation(inv2, nil)
		assert.ErrorIs(t, err, ErrUserIsAlreadyInvited)
	})

	t.Run("успешное добавление приглашения", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.New()
		chat, err := NewChat("chatName", chiefID, nil)
		require.NoError(t, err)

		// Создать и добавить приглашение
		inv, err := NewInvitation(chiefID, uuid.New())
		require.NoError(t, err)
		err = chat.AddInvitation(inv, nil)
		assert.NoError(t, err)
		assert.Contains(t, chat.Invitations, inv)
	})

	t.Run("после завершения операции, будут созданы события", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.New()
		chat, _ := NewChat("chatName", chiefID, nil)

		// Создать приглашение
		inv, _ := NewInvitation(chiefID, uuid.New())

		// Инициализировать буфер событий
		eventsBuf := new(events.Buffer)

		// Добавить приглашение
		err := chat.AddInvitation(inv, eventsBuf)
		assert.NoError(t, err)

		// Проверить, что события созданы
		require.Len(t, eventsBuf.Events(), 1)
		invitationAdded := eventsBuf.Events()[0]
		// Содержит нужных получателей
		assert.Contains(t, invitationAdded.Recipients, chat.ChiefID)
		assert.Contains(t, invitationAdded.Recipients, inv.RecipientID)
		assert.Contains(t, invitationAdded.Recipients, inv.SubjectID)
		// Содержит нужное приглашение
		assert.Equal(t, inv, invitationAdded.Data["invitation"].(Invitation))
	})
}

// TestChat_RemoveInvitation тестирует удаление приглашения из чата.
func TestChat_RemoveInvitation(t *testing.T) {
	t.Run("нельзя удалить несуществующее приглашение", func(t *testing.T) {
		// Создать чат
		chat, err := NewChat("chatName", uuid.New(), nil)
		require.NoError(t, err)

		// Удалить приглашение
		err = chat.RemoveInvitation(uuid.New(), nil)
		assert.ErrorIs(t, err, ErrInvitationNotExists)
	})

	t.Run("успешное удаление приглашения", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.New()
		chat, _ := NewChat("chatName", chiefID, nil)

		// Создать и добавить приглашение
		inv, _ := NewInvitation(chiefID, uuid.New())
		_ = chat.AddInvitation(inv, nil)

		// Удалить приглашение
		err := chat.RemoveInvitation(inv.ID, nil)
		assert.NoError(t, err)
		assert.NotContains(t, chat.Invitations, inv)
	})

	t.Run("после завершения операции, будут созданы события", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.New()
		chat, _ := NewChat("chatName", chiefID, nil)

		// Создать и добавить приглашение
		inv, _ := NewInvitation(chiefID, uuid.New())
		_ = chat.AddInvitation(inv, nil)

		// Инициализировать буфер событий
		eventsBuf := new(events.Buffer)

		// Удалить приглашение
		err := chat.RemoveInvitation(inv.ID, eventsBuf)
		assert.NoError(t, err)
		assert.NotContains(t, chat.Invitations, inv)

		// Проверить, что события созданы
		require.Len(t, eventsBuf.Events(), 1)
		invitationRemoved := eventsBuf.Events()[0]
		// Содержит нужных получателей
		assert.Contains(t, invitationRemoved.Recipients, chat.ChiefID)
		assert.Contains(t, invitationRemoved.Recipients, inv.RecipientID)
		assert.Contains(t, invitationRemoved.Recipients, inv.SubjectID)
		// Содержит нужное приглашение
		assert.Equal(t, inv, invitationRemoved.Data["invitation"].(Invitation))
	})
}

// TestChat_InvitationMethods тестирует методы для работы с приглашениями в чате.
func TestChat_InvitationMethods(t *testing.T) {
	t.Run("проверка наличия приглашения", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.New()
		chat, _ := NewChat("chatName", chiefID, nil)

		// Создать и добавить приглашение
		inv, _ := NewInvitation(chiefID, uuid.New())
		_ = chat.AddInvitation(inv, nil)

		// Проверка наличия приглашения по ID
		assert.True(t, chat.HasInvitation(inv.ID))
		assert.False(t, chat.HasInvitation(uuid.New()))
	})

	t.Run("получение приглашения по ID", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.New()
		chat, _ := NewChat("chatName", chiefID, nil)

		// Создать и добавить приглашение
		inv, _ := NewInvitation(chiefID, uuid.New())
		_ = chat.AddInvitation(inv, nil)

		// Получение приглашения по ID
		found, err := chat.Invitation(inv.ID)
		assert.NoError(t, err)
		assert.Equal(t, inv, found)
	})

	t.Run("проверка наличия приглашения по получателю", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.New()
		recipientID := uuid.New()
		chat, _ := NewChat("chatName", chiefID, nil)

		// Создать и добавить приглашение
		inv, _ := NewInvitation(chiefID, recipientID)
		_ = chat.AddInvitation(inv, nil)

		// Проверка наличия приглашения по получателю
		assert.True(t, chat.HasInvitationWithRecipient(recipientID))
		assert.False(t, chat.HasInvitationWithRecipient(uuid.New()))
	})
}
