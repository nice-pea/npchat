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
	suite.rr.users = suite.factory.NewUserrRepository()
	suite.rr.sessions = suite.factory.NewSessionnRepository()

	// Инициализация адаптеров
	suite.mockOAuthUsers = suite.generateMockUsers()
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
func (suite *servicesTestSuite) randomOAuthToken() userr.OpenAuthToken {
	t, err := userr.NewOpenAuthToken(
		randomString(32), // AccessToken
		"Bearer",         // TokenType
		randomString(32), // RefreshToken
		time.Now().Add(time.Hour*24*time.Duration(rand.Intn(7)+1)), // Expiry

	)
	suite.Require().NoError(err)

	return t
}

// Генерация случайного OAuthUser
func (suite *servicesTestSuite) randomOAuthUser() userr.OpenAuthUser {
	u, err := userr.NewOpenAuthUser(
		randomString(21),                                      // ID
		(&oauthProvider.Mock{}).Name(),                        // Provider
		randomString(8)+"@example.com",                        // Email
		randomString(6)+" "+randomString(7),                   // Name
		"https://example.com/avatar/"+randomString(10)+".png", // Picture
		userr.OpenAuthToken{},                                 // Token
	)
	suite.Require().NoError(err)

	return u
}

func (suite *servicesTestSuite) generateMockUsers() map[userr.OpenAuthToken]userr.OpenAuthUser {
	tokenToUser := make(map[userr.OpenAuthToken]userr.OpenAuthUser)

	for i := 0; i < 10; i++ {
		token := suite.randomOAuthToken()
		user := suite.randomOAuthUser()
		tokenToUser[token] = user
	}

	return tokenToUser
}
