package service

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/saime-0/nice-pea-chat/internal/adapter/oauthProvider"
	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
	"github.com/saime-0/nice-pea-chat/internal/domain/sessionn"
	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
)

type servicesTestSuite struct {
	suite.Suite
	factory *sqlite.RepositoryFactory
	rr      struct {
		chats    chatt.Repository
		sessions sessionn.Repository
		users    userr.Repository
	}
	ss struct {
		chats    *Chats
		sessions *Sessions
		users    *Users
	}
	ad struct {
		oauth OAuthProvider
	}
	mockOAuthTokens map[string]userr.OpenAuthToken
	mockOAuthUsers  map[userr.OpenAuthToken]userr.OpenAuthUser
}

func Test_ServicesTestSuite(t *testing.T) {
	suite.Run(t, new(servicesTestSuite))
}

// SetupSubTest выполняется перед каждым подтестом, связанным с suite
func (suite *servicesTestSuite) SetupSubTest() {
	var err error
	require := suite.Require()

	// Инициализация SQLiteMemory
	suite.factory, err = sqlite.InitRepositoryFactory(sqlite.Config{
		MigrationsDir: "../../migrations/repository/sqlite",
	})
	require.NoError(err)

	// Инициализация репозиториев
	suite.rr.chats = suite.factory.NewChatsRepository()
	suite.rr.users = suite.factory.NewUsersRepository()
	suite.rr.sessions = suite.factory.NewSessionsRepository()

	// Инициализация адаптеров
	suite.mockOAuthUsers = generateMockUsers()
	suite.mockOAuthTokens = make(map[string]userr.OpenAuthToken, len(suite.mockOAuthUsers))
	for token := range suite.mockOAuthUsers {
		suite.mockOAuthTokens[randomString(13)] = token
	}
	suite.ad.oauth = &oauthProvider.Mock{
		ExchangeFunc: func(code string) (userr.OpenAuthToken, error) {
			token, ok := suite.mockOAuthTokens[code]
			if !ok {
				return userr.OpenAuthToken{}, errors.New("token not found")
			}
			return token, nil
		},
		UserFunc: func(token userr.OpenAuthToken) (userr.OpenAuthUser, error) {
			user, ok := suite.mockOAuthUsers[token]
			if !ok {
				return userr.OpenAuthUser{}, errors.New("user not found")
			}
			return user, nil
		},
		AuthorizationURLFunc: func(state string) string {
			return "https://provider.com/o/oauth2/auth?code=somecode&state=" + state
		},
	}

	// Создание сервисов
	suite.ss.chats = &Chats{
		Repo: suite.rr.chats,
	}
	suite.ss.sessions = &Sessions{
		Repo: suite.rr.sessions,
	}
	suite.ss.users = &Users{
		Providers:    OAuthProviders{suite.ad.oauth.Name(): suite.ad.oauth},
		Repo:         suite.rr.users,
		SessionsRepo: suite.rr.sessions,
	}
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *servicesTestSuite) TearDownSubTest() {
	err := suite.factory.Close()
	suite.Require().NoError(err)
}

// upsertChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) upsertChat(chat chatt.Chat) chatt.Chat {
	err := suite.rr.chats.Upsert(chat)
	suite.Require().NoError(err)

	return chat
}

// upsertChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) rndChat() chatt.Chat {
	chat, err := chatt.NewChat(gofakeit.Noun(), uuid.NewString())
	suite.Require().NoError(err)

	return chat
}

// rndParticipant создает случайного участника
func (suite *servicesTestSuite) rndParticipant() chatt.Participant {
	p, err := chatt.NewParticipant(uuid.NewString())
	suite.Require().NoError(err)
	return p
}

// rndParticipant создает случайного участника
func (suite *servicesTestSuite) newParticipant(userID string) chatt.Participant {
	p, err := chatt.NewParticipant(userID)
	suite.Require().NoError(err)
	return p
}

func (suite *servicesTestSuite) addRndParticipant(chat *chatt.Chat) chatt.Participant {
	p, err := chatt.NewParticipant(uuid.NewString())
	suite.Require().NoError(err)
	suite.Require().NoError(chat.AddParticipant(p))

	return p
}

func (suite *servicesTestSuite) addParticipant(chat *chatt.Chat, p chatt.Participant) {
	suite.Require().NoError(chat.AddParticipant(p))
}

func (suite *servicesTestSuite) newInvitation(subjectID, recipientID string) chatt.Invitation {
	i, err := chatt.NewInvitation(subjectID, recipientID)
	suite.Require().NoError(err)
	return i
}

func (suite *servicesTestSuite) addInvitation(chat *chatt.Chat, i chatt.Invitation) {
	suite.Require().NoError(chat.AddInvitation(i))
}

//
//// saveMember сохраняет участника в репозиторий, в случае ошибки завершит тест
//func (suite *servicesTestSuite) saveMember(participant chatt.Participant) chatt.Participant {
//	err := suite.rr..Save(participant)
//	suite.Require().NoError(err)
//
//	return participant
//}

//// saveInvitation сохраняет приглашение в репозиторий, в случае ошибки завершит тест
//func (suite *servicesTestSuite) saveInvitation(invitation domain.Invitation) domain.Invitation {
//	err := suite.rr.invitations.Save(invitation)
//	suite.Require().NoError(err)
//
//	return invitation
//}

// saveUser сохраняет пользователя в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveUser(user userr.User) userr.User {
	err := suite.rr.users.Upsert(user)
	suite.Require().NoError(err)

	return user
}

// Функция для генерации случайной строки
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Генерация случайного OAuthToken
func randomOAuthToken() userr.OpenAuthToken {
	return userr.OpenAuthToken{
		AccessToken:  randomString(32),
		TokenType:    "Bearer",
		RefreshToken: randomString(32),
		Expiry:       time.Now().Add(time.Hour * 24 * time.Duration(rand.Intn(7)+1)),
	}
}

// Генерация случайного OAuthUser
func randomOAuthUser() userr.OpenAuthUser {
	return userr.OpenAuthUser{
		ID:      randomString(21), // Providers ID обычно длина ~21
		Email:   randomString(8) + "@example.com",
		Name:    randomString(6) + " " + randomString(7),
		Picture: "https://example.com/avatar/" + randomString(10) + ".png",
	}
}

func generateMockUsers() map[userr.OpenAuthToken]userr.OpenAuthUser {
	tokenToUser := make(map[userr.OpenAuthToken]userr.OpenAuthUser)

	for i := 0; i < 10; i++ {
		token := randomOAuthToken()
		user := randomOAuthUser()
		tokenToUser[token] = user
	}

	return tokenToUser
}

func randomPassword() string {
	validPasswords := []string{
		"Ab1!xyzZ",
		"Password123!",
		"Пароль123!",
		"Aa1~!?@#$%^&*_-+()[]{}></\\|\"'.,:;",
		"P@ssw0rd_123",
		"Passворд123!",
	}

	return validPasswords[rand.Intn(len(validPasswords))]
}
