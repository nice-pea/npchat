package chatt

import (
	"math/rand"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
)

// НАСТРОЙКА ТЕСТОВОГО ОКРУЖЕНИЯ

type repoSuite struct {
	suite.Suite
	newRepository func() Repository
	repo          Repository
}

// Test_repoSuite запускает набор тестов
func Test_repoSuite(t *testing.T) {
	suite.Run(t, &repoSuite{
		Suite:         suite.Suite{},
		newRepository: sqliteRepoConstructor(t),
		repo:          nil,
	})
}

// SetupSubTest подготавливает репозиторий для каждого подтеста
func (suite *repoSuite) SetupSubTest() {
	suite.repo = suite.newRepository()
}

func (suite *repoSuite) TearDownSubTest() {}

// sqliteRepoConstructor создает конструктор SQLite репозитория с тестовой конфигурацией
func sqliteRepoConstructor(t *testing.T) func() Repository {
	sqliteConfig := sqlite.Config{
		MigrationsDir: "../../../migrations/repository/sqlite",
	}
	repositoryFactory, err := sqlite.InitRepositoryFactory(sqliteConfig)
	assert.Nil(t, err)
	assert.NotNil(t, repositoryFactory)
	repo := repositoryFactory.NewChatsRepository()
	require.NotNil(t, repo)

	f, err := sqlite.InitRepositoryFactory(sqliteConfig)
	require.NoError(t, err)
	return f.NewChattRepository
}

// rndChat создает случайный экземпляр чата
func (suite *repoSuite) rndChat() Chat {
	chat, err := NewChat(gofakeit.Noun(), uuid.NewString())
	suite.Require().NoError(err)

	return chat
}

// upsertChat сохраняет чат в репозиторий
func (suite *repoSuite) upsertChat(chat Chat) Chat {
	err := suite.repo.Upsert(chat)
	suite.Require().NoError(err)

	return chat
}

// rndElem возвращает случайный элемент из среза
func rndElem[T any](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	index := rand.Intn(len(slice))
	return slice[index]
}

// rndInv создает случайное приглашение
func (suite *repoSuite) rndInv() Invitation {
	inv, err := NewInvitation(uuid.NewString(), uuid.NewString())
	suite.Require().NoError(err)
	return inv
}

// rndParticipant создает случайного участника
func (suite *repoSuite) rndParticipant() Participant {
	p, err := NewParticipant(uuid.NewString())
	suite.Require().NoError(err)
	return p
}

// addRndParticipant добавляет случайного участника в чат
func (suite *repoSuite) addRndParticipant(chat *Chat) {
	p, err := NewParticipant(uuid.NewString())
	suite.Require().NoError(err)
	suite.Require().NoError(chat.AddParticipant(p))
}

// addRndInv добавляет случайное приглашение в чат
func (suite *repoSuite) addRndInv(chat *Chat) {
	inv, err := NewInvitation(rndElem(chat.Participants).UserID, uuid.NewString())
	suite.Require().NoError(err)
	suite.Require().NoError(chat.AddInvitation(inv))
}

// ТЕСТЫ

// Test_Repository реализацию репозитория
func (suite *repoSuite) Test_Repository() {
	suite.Run("List", func() {
		suite.Run("из пустого репозитория вернется пустой список", func() {
			chats, err := suite.repo.List(Filter{})
			suite.NoError(err)
			suite.Empty(chats)
		})

		suite.Run("без фильтра из репозитория вернутся все сохраненные чаты", func() {
			chats := make([]Chat, 10)
			for i := range chats {
				chats[i] = suite.upsertChat(suite.rndChat())
			}
			chatsFromRepo, err := suite.repo.List(Filter{})
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
			chatsFromRepo, err := suite.repo.List(Filter{
				ID: expectedChat.ID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, 1)
			suite.Equal(expectedChat, chatsFromRepo[0])
		})

		suite.Run("с фильтром по InvitationID вернутся чаты, имеющие с приглашение с таким ID", func() {
			// Создать много чатов
			chats := make([]Chat, 10)
			for i := range chats {
				chats[i] = suite.rndChat()
				suite.addRndParticipant(&chats[i])
				suite.addRndInv(&chats[i])
				suite.upsertChat(chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := rndElem(chats)

			// Получить список
			chatsFromRepo, err := suite.repo.List(Filter{
				InvitationID: expectedChat.Invitations[0].ID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, 1)
			suite.Equal(expectedChat, chatsFromRepo[0])
		})

		suite.Run("с фильтром по InvitationRecipientID вернутся чаты, имеющие с приглашения, направленные пользователю с тем ID", func() {
			// Создать много чатов
			chats := make([]Chat, 10)
			for i := range chats {
				chats[i] = suite.rndChat()
				suite.addRndParticipant(&chats[i])
				suite.addRndInv(&chats[i])
				suite.upsertChat(chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := rndElem(chats)

			// Получить список
			chatsFromRepo, err := suite.repo.List(Filter{
				InvitationRecipientID: expectedChat.Invitations[0].RecipientID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, 1)
			suite.Equal(expectedChat, chatsFromRepo[0])
		})

		suite.Run("с фильтром по ParticipantID вернутся чаты, в которых состоит пользователь с тем ID", func() {
			// Создать много чатов
			chats := make([]Chat, 10)
			for i := range chats {
				chats[i] = suite.rndChat()
				suite.addRndParticipant(&chats[i])
				suite.upsertChat(chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := rndElem(chats)

			// Получить список
			chatsFromRepo, err := suite.repo.List(Filter{
				ParticipantID: expectedChat.Participants[0].UserID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, 1)
			suite.Equal(expectedChat, chatsFromRepo[0])
		})

		suite.Run("можно искать по всем фильтрам сразу", func() {
			// Создать много чатов
			chats := make([]Chat, 10)
			for i := range chats {
				chats[i] = suite.rndChat()
				suite.addRndParticipant(&chats[i])
				suite.addRndInv(&chats[i])
				suite.upsertChat(chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := rndElem(chats)

			// Получить список
			chatsFromRepo, err := suite.repo.List(Filter{
				ID:                    expectedChat.ID,
				InvitationID:          expectedChat.Invitations[0].ID,
				InvitationRecipientID: expectedChat.Invitations[0].RecipientID,
				ParticipantID:         expectedChat.Participants[0].UserID,
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
				err := chat.AddParticipant(rndp)
				suite.Require().NoError(err)
				suite.upsertChat(chat)
			}
			// Создать несколько других чатов
			for range 21 {
				suite.upsertChat(suite.rndChat())
			}
			// Получить список
			chatsFromRepo, err := suite.repo.List(Filter{
				ParticipantID: rndp.UserID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(chatsFromRepo, expectedCount)
		})
	})

	suite.Run("Upsert", func() {
		suite.Run("нельзя сохранять чат без ID", func() {
			err := suite.repo.Upsert(Chat{
				ID:   "",
				Name: "someName",
			})
			suite.Error(err)
		})

		suite.Run("остальные поля, кроме ID могут быть пустыми", func() {
			err := suite.repo.Upsert(Chat{
				ID: uuid.NewString(),
			})
			suite.NoError(err)
		})

		suite.Run("сохраненный чат полностью соответствует сохраняемому", func() {
			// Наполнить чат
			chat := suite.rndChat()
			suite.addRndParticipant(&chat)
			suite.addRndInv(&chat)

			// Сохранить чат
			err := suite.repo.Upsert(chat)
			suite.Require().NoError(err)

			// Прочитать из репозитория
			chats, err := suite.repo.List(Filter{})
			suite.NoError(err)
			suite.Require().Len(chats, 1)
			suite.Equal(chat, chats[0])
		})

		suite.Run("перезапись с новыми значениями по ID", func() {
			id := uuid.NewString()
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
			chats, err := suite.repo.List(Filter{})
			suite.NoError(err)
			suite.Require().Len(chats, 1)
			suite.Equal(expectedChat, chats[0])
		})
	})
}
