package pgsqlRepository

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/common"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

func (suite *Suite) Test_ChattRepository() {
	suite.Run("List", func() {
		suite.Run("из пустого репозитория вернется пустой список", func() {
			chats, err := suite.RR.Chats.List(chatt.Filter{})
			suite.NoError(err)
			suite.Empty(chats)
		})

		suite.Run("без фильтра из репозитория вернутся все сохраненные чаты", func() {
			chats := make([]chatt.Chat, 10)
			for i := range chats {
				chats[i] = suite.upsertChat(suite.rndChat())
			}
			chatsFromRepo, err := suite.RR.Chats.List(chatt.Filter{})
			suite.NoError(err)
			suite.Len(chatsFromRepo, len(chats))
		})

		suite.Run("с фильтром по ID вернется сохраненный чат", func() {
			// Создать много чатов
			for range 10 {
				suite.upsertChat(suite.rndChat())
			}
			// Определить случайны искомый чат
			expectedChat := suite.upsertChat(suite.rndChat())
			// Получить список
			chatsFromRepo, err := suite.RR.Chats.List(chatt.Filter{
				ID: expectedChat.ID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, 1)
			suite.Equal(expectedChat, chatsFromRepo[0])
		})

		suite.Run("с фильтром по InvitationID вернутся чаты, имеющие с приглашение с таким ID", func() {
			// Создать много чатов
			chats := make([]chatt.Chat, 10)
			for i := range chats {
				chats[i] = suite.rndChat()
				suite.addRndParticipant(&chats[i])
				suite.addRndInv(&chats[i])
				suite.upsertChat(chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := common.RndElem(chats)

			// Получить список
			chatsFromRepo, err := suite.RR.Chats.List(chatt.Filter{
				InvitationID: expectedChat.Invitations[0].ID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, 1)
			suite.Equal(expectedChat, chatsFromRepo[0])
		})

		suite.Run("с фильтром по InvitationRecipientID вернутся чаты, имеющие с приглашения, направленные пользователю с тем ID", func() {
			// Создать много чатов
			chats := make([]chatt.Chat, 10)
			for i := range chats {
				chats[i] = suite.rndChat()
				suite.addRndParticipant(&chats[i])
				suite.addRndInv(&chats[i])
				suite.upsertChat(chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := common.RndElem(chats)

			// Получить список
			chatsFromRepo, err := suite.RR.Chats.List(chatt.Filter{
				InvitationRecipientID: expectedChat.Invitations[0].RecipientID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, 1)
			suite.Equal(expectedChat, chatsFromRepo[0])
		})

		suite.Run("с фильтром по ParticipantID вернутся чаты, в которых состоит пользователь с тем ID", func() {
			// Создать много чатов
			chats := make([]chatt.Chat, 10)
			for i := range chats {
				chats[i] = suite.rndChat()
				suite.addRndParticipant(&chats[i])
				suite.upsertChat(chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := common.RndElem(chats)

			// Получить список
			chatsFromRepo, err := suite.RR.Chats.List(chatt.Filter{
				ParticipantID: expectedChat.Participants[0].UserID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, 1)
			suite.Equal(expectedChat, chatsFromRepo[0])
		})

		suite.Run("с фильтром ActiveBefore вернутся чаты с меньшим LastActiveAt", func() {
			// Создать чаты с разными LastActiveAt
			now := time.Now().Truncate(time.Microsecond)
			for _, duration := range []time.Duration{
				time.Hour,
				time.Second,
				time.Hour * 2,
				time.Minute,
			} {
				chat := suite.rndChat()
				_ = chat.SetLastActiveAt(now.Add(duration), nil)
				suite.upsertChat(chat)
			}

			// Получить список
			chatsFromRepo, err := suite.RR.Chats.List(chatt.Filter{
				ActiveBefore: now.Add(time.Hour),
			})
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, 2)
			suite.True(now.Add(time.Minute).Equal(chatsFromRepo[0].LastActiveAt))
			suite.True(now.Add(time.Second).Equal(chatsFromRepo[1].LastActiveAt))
		})

		suite.Run("с limit вернется ограниченное количество элементов", func() {
			// Создать чаты
			const limit = 10
			var createdChats []chatt.Chat
			for range limit * 2 {
				createdChats = append(createdChats, suite.upsertChat(suite.rndChat()))
				time.Sleep(time.Millisecond * 10)
			}

			// Получить список
			chatsFromRepo, err := suite.RR.Chats.List(chatt.Filter{
				Limit: limit,
			})
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, limit)

			for i := 0; i < limit; i++ {
				expectedIdx := len(createdChats) - 1 - i
				suite.Equal(createdChats[expectedIdx], chatsFromRepo[i])
			}
		})

		suite.Run("можно искать по всем фильтрам сразу", func() {
			// Создать много чатов
			chats := make([]chatt.Chat, 10)
			for i := range chats {
				chats[i] = suite.rndChat()
				suite.addRndParticipant(&chats[i])
				suite.addRndInv(&chats[i])
				suite.upsertChat(chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := common.RndElem(chats)

			// Получить список
			chatsFromRepo, err := suite.RR.Chats.List(chatt.Filter{
				ID:                    expectedChat.ID,
				InvitationID:          expectedChat.Invitations[0].ID,
				InvitationRecipientID: expectedChat.Invitations[0].RecipientID,
				ParticipantID:         expectedChat.Participants[0].UserID,
				ActiveBefore:          time.Now().Add(time.Second),
				Limit:                 len(chats),
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, 1)
			suite.Equal(expectedChat, chatsFromRepo[0])
		})

		suite.Run("можно вернуться несколько элементов", func() {
			// Участник, который есть во многих чатах
			rndp := suite.rndParticipant()
			// Создать много чатов с искомым участником
			const expectedCount = 10
			for range expectedCount {
				chat := suite.rndChat()
				err := chat.AddParticipant(rndp, nil)
				suite.Require().NoError(err)
				suite.upsertChat(chat)
			}
			// Создать несколько других чатов
			for range 21 {
				suite.upsertChat(suite.rndChat())
			}
			// Получить список
			chatsFromRepo, err := suite.RR.Chats.List(chatt.Filter{
				ParticipantID: rndp.UserID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, expectedCount)
		})

		suite.Run("чаты возвращаются в порядке убывания LastActiveAt", func() {
			// Создать чаты с разными LastActiveAt
			for _, duration := range []time.Duration{
				time.Hour,
				time.Second,
				time.Hour * 2,
				time.Minute,
			} {
				chat := suite.rndChat()
				_ = chat.SetLastActiveAt(time.Now().Add(duration), nil)
				suite.upsertChat(chat)
			}

			// Получить список
			chatsFromRepo, err := suite.RR.Chats.List(chatt.Filter{})
			// Сравнить ожидания и результат
			suite.NoError(err)

			prevActiveAt := chatsFromRepo[0].LastActiveAt
			for _, chat := range chatsFromRepo[1:] {
				suite.Less(chat.LastActiveAt, prevActiveAt)
				prevActiveAt = chat.LastActiveAt
			}
		})
	})

	suite.Run("Upsert", func() {
		suite.Run("нельзя сохранять чат без ID", func() {
			err := suite.RR.Chats.Upsert(chatt.Chat{
				ID:   uuid.Nil,
				Name: "someName",
			})
			suite.Error(err)
		})

		suite.Run("остальные поля, кроме ID могут быть пустыми", func() {
			err := suite.RR.Chats.Upsert(chatt.Chat{
				ID: uuid.New(),
			})
			suite.NoError(err)
		})

		suite.Run("сохраненный чат полностью соответствует сохраняемому", func() {
			// Наполнить чат
			chat := suite.rndChat()
			suite.addRndParticipant(&chat)
			suite.addRndInv(&chat)

			// Сохранить чат
			err := suite.RR.Chats.Upsert(chat)
			suite.Require().NoError(err)

			// Прочитать из репозитория
			chats, err := suite.RR.Chats.List(chatt.Filter{})
			suite.NoError(err)
			suite.Require().Len(chats, 1)
			suite.Equal(chat, chats[0])
		})

		suite.Run("перезапись с новыми значениями по ID", func() {
			id := uuid.New()
			// Несколько промежуточных состояний чата
			for range 33 {
				chat := suite.rndChat()
				chat.ID = id
				suite.upsertChat(chat)
			}
			// Последнее сохраненное состояние чата
			expectedChat := suite.rndChat()
			expectedChat.ID = id
			suite.upsertChat(expectedChat)

			// Прочитать из репозитория
			chats, err := suite.RR.Chats.List(chatt.Filter{})
			suite.NoError(err)
			suite.Require().Len(chats, 1)
			suite.Equal(expectedChat, chats[0])
		})
	})
}

// rndChat создает случайный экземпляр чата
func (suite *Suite) rndChat() chatt.Chat {
	suite.T().Helper()
	chat, err := chatt.NewChat(gofakeit.Noun(), uuid.New(), nil)
	suite.Require().NoError(err)

	return chat
}

// upsertChat сохраняет чат в репозиторий
func (suite *Suite) upsertChat(chat chatt.Chat) chatt.Chat {
	suite.T().Helper()
	err := suite.RR.Chats.Upsert(chat)
	suite.Require().NoError(err)

	return chat
}

// rndParticipant создает случайного участника
func (suite *Suite) rndParticipant() chatt.Participant {
	suite.T().Helper()
	p, err := chatt.NewParticipant(uuid.New())
	suite.Require().NoError(err)
	return p
}

// addRndParticipant добавляет случайного участника в чат
func (suite *Suite) addRndParticipant(chat *chatt.Chat) {
	suite.T().Helper()
	p, err := chatt.NewParticipant(uuid.New())
	suite.Require().NoError(err)
	suite.Require().NoError(chat.AddParticipant(p, nil))
}

// addRndInv добавляет случайное приглашение в чат
func (suite *Suite) addRndInv(chat *chatt.Chat) {
	suite.T().Helper()
	inv, err := chatt.NewInvitation(common.RndElem(chat.Participants).UserID, uuid.New())
	suite.Require().NoError(err)
	suite.Require().NoError(chat.AddInvitation(inv, nil))
}
