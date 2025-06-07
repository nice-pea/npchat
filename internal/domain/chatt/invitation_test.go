package chatt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewInvitation тестирует создание приглашения.
func TestNewInvitation(t *testing.T) {
	t.Run("параметр subjectID не должен быть пустым", func(t *testing.T) {
		inv, err := NewInvitation("", uuid.NewString())
		assert.Zero(t, inv)
		assert.Error(t, err)
	})

	t.Run("параметр recipientID не должен быть пустым", func(t *testing.T) {
		inv, err := NewInvitation(uuid.NewString(), "")
		assert.Zero(t, inv)
		assert.Error(t, err)
	})

	t.Run("параметры должны быть валидными UUID", func(t *testing.T) {
		inv, err := NewInvitation("invalid", "invalid")
		assert.Zero(t, inv)
		assert.Error(t, err)
	})

	t.Run("subjectID и recipientID не могут быть одинаковыми", func(t *testing.T) {
		id := uuid.NewString()
		inv, err := NewInvitation(id, id)
		assert.Zero(t, inv)
		assert.ErrorIs(t, err, ErrSubjectAndRecipientMustBeDifferent)
	})

	t.Run("создание валидного приглашения", func(t *testing.T) {
		subjectID := uuid.NewString()
		recipientID := uuid.NewString()
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
		chat, err := NewChat("chatName", uuid.NewString())
		require.NoError(t, err)

		// Создать и добавить первое приглашение
		inv, err := NewInvitation(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)
		err = chat.AddInvitation(inv)
		assert.ErrorIs(t, err, ErrSubjectIsNotMember)
	})

	t.Run("нельзя пригласить существующего участника", func(t *testing.T) {
		// Создать чат
		chief := uuid.NewString()
		chat, err := NewChat("test", chief)
		require.NoError(t, err)

		// Создать и добавить участника
		p, err := NewParticipant(uuid.NewString())
		require.NoError(t, err)
		err = chat.AddParticipant(p)
		require.NoError(t, err)

		// Добавить приглашение пользователю, который уже участник
		inv, err := NewInvitation(chief, p.UserID)
		require.NoError(t, err)
		err = chat.AddInvitation(inv)
		assert.ErrorIs(t, err, ErrParticipantExists)
	})

	t.Run("нельзя пригласить уже приглашенного пользователя", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.NewString()
		recipientID := uuid.NewString()
		chat, err := NewChat("chatName", chiefID)
		require.NoError(t, err)

		// Создать и добавить первое приглашение
		inv1, err := NewInvitation(chiefID, recipientID)
		require.NoError(t, err)
		err = chat.AddInvitation(inv1)
		require.NoError(t, err)

		//  Создать и добавить второе приглашение
		inv2, err := NewInvitation(chiefID, recipientID)
		require.NoError(t, err)
		err = chat.AddInvitation(inv2)
		assert.ErrorIs(t, err, ErrUserIsAlreadyInvited)
	})

	t.Run("успешное добавление приглашения", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.NewString()
		chat, err := NewChat("chatName", chiefID)
		require.NoError(t, err)

		// Создать и добавить приглашение
		inv, err := NewInvitation(chiefID, uuid.NewString())
		require.NoError(t, err)
		err = chat.AddInvitation(inv)
		assert.NoError(t, err)
		assert.Contains(t, chat.Invitations, inv)
	})
}

// TestChat_RemoveInvitation тестирует удаление приглашения из чата.
func TestChat_RemoveInvitation(t *testing.T) {
	t.Run("нельзя удалить несуществующее приглашение", func(t *testing.T) {
		// Создать чат
		chat, err := NewChat("chatName", uuid.NewString())
		require.NoError(t, err)

		// Удалить приглашение
		err = chat.RemoveInvitation(uuid.NewString())
		assert.ErrorIs(t, err, ErrInvitationNotExists)
	})

	t.Run("успешное удаление приглашения", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.NewString()
		chat, _ := NewChat("chatName", chiefID)

		// Создать и добавить приглашение
		inv, _ := NewInvitation(chiefID, uuid.NewString())
		_ = chat.AddInvitation(inv)

		// Удалить приглашение
		err := chat.RemoveInvitation(inv.ID)
		assert.NoError(t, err)
		assert.NotContains(t, chat.Invitations, inv)
	})
}

// TestChat_InvitationMethods тестирует методы для работы с приглашениями в чате.
func TestChat_InvitationMethods(t *testing.T) {
	t.Run("проверка наличия приглашения", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.NewString()
		chat, _ := NewChat("chatName", chiefID)

		// Создать и добавить приглашение
		inv, _ := NewInvitation(chiefID, uuid.NewString())
		_ = chat.AddInvitation(inv)

		// Проверка наличия приглашения по ID
		assert.True(t, chat.HasInvitation(inv.ID))
		assert.False(t, chat.HasInvitation(uuid.NewString()))
	})

	t.Run("получение приглашения по ID", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.NewString()
		chat, _ := NewChat("chatName", chiefID)

		// Создать и добавить приглашение
		inv, _ := NewInvitation(chiefID, uuid.NewString())
		_ = chat.AddInvitation(inv)

		// Получение приглашения по ID
		found, err := chat.Invitation(inv.ID)
		assert.NoError(t, err)
		assert.Equal(t, inv, found)
	})

	t.Run("проверка наличия приглашения по получателю", func(t *testing.T) {
		// Создать чат
		chiefID := uuid.NewString()
		recipientID := uuid.NewString()
		chat, _ := NewChat("chatName", chiefID)

		// Создать и добавить приглашение
		inv, _ := NewInvitation(chiefID, recipientID)
		_ = chat.AddInvitation(inv)

		// Проверка наличия приглашения по получателю
		assert.True(t, chat.HasInvitationWithRecipient(recipientID))
		assert.False(t, chat.HasInvitationWithRecipient(uuid.NewString()))
	})
}
