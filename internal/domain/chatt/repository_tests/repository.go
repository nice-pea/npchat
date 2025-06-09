package repository_tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/common"
	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
)

// TestRepository реализацию репозитория
func TestRepository(t *testing.T, newRepository func() chatt.Repository) {
	t.Run("List", func(t *testing.T) {
		t.Run("из пустого репозитория вернется пустой список", func(t *testing.T) {
			r := newRepository()
			chats, err := r.List(chatt.Filter{})
			assert.NoError(t, err)
			assert.Empty(t, chats)
		})

		t.Run("без фильтра из репозитория вернутся все сохраненные чаты", func(t *testing.T) {
			r := newRepository()
			chats := make([]chatt.Chat, 10)
			for i := range chats {
				chats[i] = upsertChat(t, r, rndChat(t))
			}
			chatsFromRepo, err := r.List(chatt.Filter{})
			assert.NoError(t, err)
			assert.Len(t, chatsFromRepo, len(chats))
		})

		t.Run("с фильтром по ID вернется сохраненный чат", func(t *testing.T) {
			r := newRepository()
			// Создать много чатов
			for range 10 {
				upsertChat(t, r, rndChat(t))
			}
			// Определить случайны искомый чат
			expectedChat := upsertChat(t, r, rndChat(t))
			// Получить список
			chatsFromRepo, err := r.List(chatt.Filter{
				ID: expectedChat.ID,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, chatsFromRepo, 1)
			assert.Equal(t, expectedChat, chatsFromRepo[0])
		})

		t.Run("с фильтром по InvitationID вернутся чаты, имеющие с приглашение с таким ID", func(t *testing.T) {
			r := newRepository()
			// Создать много чатов
			chats := make([]chatt.Chat, 10)
			for i := range chats {
				chats[i] = rndChat(t)
				addRndParticipant(t, &chats[i])
				addRndInv(t, &chats[i])
				upsertChat(t, r, chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := common.RndElem(chats)

			// Получить список
			chatsFromRepo, err := r.List(chatt.Filter{
				InvitationID: expectedChat.Invitations[0].ID,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, chatsFromRepo, 1)
			assert.Equal(t, expectedChat, chatsFromRepo[0])
		})

		t.Run("с фильтром по InvitationRecipientID вернутся чаты, имеющие с приглашения, направленные пользователю с тем ID", func(t *testing.T) {
			r := newRepository()
			// Создать много чатов
			chats := make([]chatt.Chat, 10)
			for i := range chats {
				chats[i] = rndChat(t)
				addRndParticipant(t, &chats[i])
				addRndInv(t, &chats[i])
				upsertChat(t, r, chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := common.RndElem(chats)

			// Получить список
			chatsFromRepo, err := r.List(chatt.Filter{
				InvitationRecipientID: expectedChat.Invitations[0].RecipientID,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, chatsFromRepo, 1)
			assert.Equal(t, expectedChat, chatsFromRepo[0])
		})

		t.Run("с фильтром по ParticipantID вернутся чаты, в которых состоит пользователь с тем ID", func(t *testing.T) {
			r := newRepository()
			// Создать много чатов
			chats := make([]chatt.Chat, 10)
			for i := range chats {
				chats[i] = rndChat(t)
				addRndParticipant(t, &chats[i])
				upsertChat(t, r, chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := common.RndElem(chats)

			// Получить список
			chatsFromRepo, err := r.List(chatt.Filter{
				ParticipantID: expectedChat.Participants[0].UserID,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, chatsFromRepo, 1)
			assert.Equal(t, expectedChat, chatsFromRepo[0])
		})

		t.Run("можно искать по всем фильтрам сразу", func(t *testing.T) {
			r := newRepository()
			// Создать много чатов
			chats := make([]chatt.Chat, 10)
			for i := range chats {
				chats[i] = rndChat(t)
				addRndParticipant(t, &chats[i])
				addRndInv(t, &chats[i])
				upsertChat(t, r, chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := common.RndElem(chats)

			// Получить список
			chatsFromRepo, err := r.List(chatt.Filter{
				ID:                    expectedChat.ID,
				InvitationID:          expectedChat.Invitations[0].ID,
				InvitationRecipientID: expectedChat.Invitations[0].RecipientID,
				ParticipantID:         expectedChat.Participants[0].UserID,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, chatsFromRepo, 1)
			assert.Equal(t, expectedChat, chatsFromRepo[0])
		})

		t.Run("можно вернуться несколько элементов", func(t *testing.T) {
			r := newRepository()
			// Участник, который есть во многих чатах
			rndp := rndParticipant(t)
			// Создать много чатов с искомым участником
			const expectedCount = 10
			for range expectedCount {
				chat := rndChat(t)
				err := chat.AddParticipant(rndp)
				require.NoError(t, err)
				upsertChat(t, r, chat)
			}
			// Создать несколько других чатов
			for range 21 {
				upsertChat(t, r, rndChat(t))
			}
			// Получить список
			chatsFromRepo, err := r.List(chatt.Filter{
				ParticipantID: rndp.UserID,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, chatsFromRepo, expectedCount)
		})
	})

	t.Run("Upsert", func(t *testing.T) {
		t.Run("нельзя сохранять чат без ID", func(t *testing.T) {
			r := newRepository()
			err := r.Upsert(chatt.Chat{
				ID:   "",
				Name: "someName",
			})
			assert.Error(t, err)
		})

		t.Run("остальные поля, кроме ID могут быть пустыми", func(t *testing.T) {
			r := newRepository()
			err := r.Upsert(chatt.Chat{
				ID: uuid.NewString(),
			})
			assert.NoError(t, err)
		})

		t.Run("сохраненный чат полностью соответствует сохраняемому", func(t *testing.T) {
			r := newRepository()
			// Наполнить чат
			chat := rndChat(t)
			addRndParticipant(t, &chat)
			addRndInv(t, &chat)

			// Сохранить чат
			err := r.Upsert(chat)
			require.NoError(t, err)

			// Прочитать из репозитория
			chats, err := r.List(chatt.Filter{})
			assert.NoError(t, err)
			require.Len(t, chats, 1)
			assert.Equal(t, chat, chats[0])
		})

		t.Run("перезапись с новыми значениями по ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			// Несколько промежуточных состояний чата
			for range 33 {
				chat := rndChat(t)
				chat.ID = id
				upsertChat(t, r, chat)
			}
			// Последнее сохраненное состояние чата
			expectedChat := rndChat(t)
			expectedChat.ID = id
			upsertChat(t, r, expectedChat)

			// Прочитать из репозитория
			chats, err := r.List(chatt.Filter{})
			assert.NoError(t, err)
			require.Len(t, chats, 1)
			assert.Equal(t, expectedChat, chats[0])
		})
	})
}

// rndChat создает случайный экземпляр чата
func rndChat(t *testing.T) chatt.Chat {
	chat, err := chatt.NewChat(gofakeit.Noun(), uuid.NewString())
	require.NoError(t, err)

	return chat
}

// upsertChat сохраняет чат в репозиторий
func upsertChat(t *testing.T, r chatt.Repository, chat chatt.Chat) chatt.Chat {
	err := r.Upsert(chat)
	require.NoError(t, err)

	return chat
}

// rndInv создает случайное приглашение
func rndInv(t *testing.T) chatt.Invitation {
	inv, err := chatt.NewInvitation(uuid.NewString(), uuid.NewString())
	require.NoError(t, err)
	return inv
}

// rndParticipant создает случайного участника
func rndParticipant(t *testing.T) chatt.Participant {
	p, err := chatt.NewParticipant(uuid.NewString())
	require.NoError(t, err)
	return p
}

// addRndParticipant добавляет случайного участника в чат
func addRndParticipant(t *testing.T, chat *chatt.Chat) {
	p, err := chatt.NewParticipant(uuid.NewString())
	require.NoError(t, err)
	require.NoError(t, chat.AddParticipant(p))
}

// addRndInv добавляет случайное приглашение в чат
func addRndInv(t *testing.T, chat *chatt.Chat) {
	inv, err := chatt.NewInvitation(common.RndElem(chat.Participants).UserID, uuid.NewString())
	require.NoError(t, err)
	require.NoError(t, chat.AddInvitation(inv))
}
