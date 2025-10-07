package serviceSuite

import (
	"errors"
	"math/rand"
	"slices"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/common"
	"github.com/nice-pea/npchat/internal/domain/chatt"
	mockChatt "github.com/nice-pea/npchat/internal/domain/chatt/mocks"
	mockSessionn "github.com/nice-pea/npchat/internal/domain/sessionn/mocks"
	mockUserr "github.com/nice-pea/npchat/internal/domain/userr/mocks"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	"github.com/nice-pea/npchat/internal/usecases/events"
	mockOauth "github.com/nice-pea/npchat/internal/usecases/users/oauth/mocks"
)

type SuiteWithMocks struct {
	testifySuite.Suite
	RR struct {
		Chats    *mockChatt.Repository
		Sessions *mockSessionn.Repository
		Users    *mockUserr.Repository
	}
	Adapters struct {
		Oauth *mockOauth.Provider
	}
	MockOauthTokens map[string]userr.OpenAuthToken
	MockOauthUsers  map[userr.OpenAuthToken]userr.OpenAuthUser
}

// SetupTest выполняется перед каждым тестом, связанным с suite
func (suite *SuiteWithMocks) SetupTest() {
	// Инициализация репозиториев
	suite.RR.Chats = mockChatt.NewRepository(suite.T())
	suite.RR.Users = mockUserr.NewRepository(suite.T())
	suite.RR.Sessions = mockSessionn.NewRepository(suite.T())

	// Инициализация адаптеров
	suite.Adapters.Oauth = mockOauth.NewProvider(suite.T())
	suite.Adapters.Oauth.
		On("Name", mock.Anything).Maybe().
		Return("mock").
		On("Exchange", mock.Anything).Maybe().
		Return(func(code string) (userr.OpenAuthToken, error) {
			token, ok := suite.MockOauthTokens[code]
			if !ok {
				return userr.OpenAuthToken{}, errors.New("token not found")
			}
			return token, nil
		}).
		On("User", mock.Anything).Maybe().
		Return(func(token userr.OpenAuthToken) (userr.OpenAuthUser, error) {
			user, ok := suite.MockOauthUsers[token]
			if !ok {
				return userr.OpenAuthUser{}, errors.New("user not found")
			}
			return user, nil
		}).
		On("AuthorizationURL", mock.Anything).Maybe().
		Return(func(state string) string {
			return "https://provider.com/o/oauth2/auth?" +
				"code=someCode" +
				"&state=" + state
		})

	// Тестовые пользователи и токены
	suite.MockOauthUsers = suite.GenerateMockUsers()
	suite.MockOauthTokens = make(map[string]userr.OpenAuthToken, len(suite.MockOauthUsers))
	for token := range suite.MockOauthUsers {
		suite.MockOauthTokens[RandomString(13)] = token
	}
}

// SetupAcceptInvitationMocks настраивает моки для успешного принятия приглашения
func (suite *SuiteWithMocks) SetupAcceptInvitationMocks(invitationID uuid.UUID, chat chatt.Chat) {
	suite.RR.Chats.EXPECT().List(chatt.Filter{
		InvitationID: invitationID,
	}).Return([]chatt.Chat{chat}, nil).Once()

	suite.RR.Chats.EXPECT().Upsert(mock.Anything).RunAndReturn(func(updatedChat chatt.Chat) error {
		suite.Equal(chat.ID, updatedChat.ID)
		return nil
	}).Once()
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *SuiteWithMocks) TearDownSubTest() {
	// пересоздаем моки репозиториев
	suite.RR.Chats = mockChatt.NewRepository(suite.T())
	suite.RR.Users = mockUserr.NewRepository(suite.T())
	suite.RR.Sessions = mockSessionn.NewRepository(suite.T())
}

// TearDownTest выполняется после каждого подтеста, связанного с suite
func (suite *SuiteWithMocks) TearDownTest() {

}

// EqualSessions сравнивает две сессии
func (suite *SuiteWithMocks) EqualSessions(s1, s2 sessionn.Session) {
	suite.Equal(s1.ID, s2.ID)
	suite.Equal(s1.UserID, s2.UserID)
	suite.Equal(s1.Name, s2.Name)
	suite.Equal(s1.Status, s2.Status)
	suite.Equal(s1.AccessToken.Token, s2.AccessToken.Token)
	suite.True(s1.AccessToken.Expiry.Equal(s2.AccessToken.Expiry))
	suite.Equal(s1.RefreshToken.Token, s2.RefreshToken.Token)
	suite.True(s1.RefreshToken.Expiry.Equal(s2.RefreshToken.Expiry))
}

// UpsertChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *SuiteWithMocks) UpsertChat(chat chatt.Chat) chatt.Chat {
	// настройка мока
	suite.RR.Chats.EXPECT().Upsert(chat).Return(nil).Maybe()
	err := suite.RR.Chats.Upsert(chat)
	suite.Require().NoError(err)

	return chat
}

// RndChat создает случайный чат
func (suite *SuiteWithMocks) RndChat() chatt.Chat {
	chat, err := chatt.NewChat(gofakeit.Noun(), uuid.New(), nil)
	suite.Require().NoError(err)

	return chat
}

// NewParticipant создает случайного участника
func (suite *SuiteWithMocks) NewParticipant(userID uuid.UUID) chatt.Participant {
	p, err := chatt.NewParticipant(userID)
	suite.Require().NoError(err)
	return p
}

// AddRndParticipant добавляет случайного участника в чат
func (suite *SuiteWithMocks) AddRndParticipant(chat *chatt.Chat) chatt.Participant {
	p, err := chatt.NewParticipant(uuid.New())
	suite.Require().NoError(err)
	suite.Require().NoError(chat.AddParticipant(p, nil))

	return p
}

// AddParticipant добавляет участника в чат
func (suite *SuiteWithMocks) AddParticipant(chat *chatt.Chat, p chatt.Participant) {
	suite.Require().NoError(chat.AddParticipant(p, nil))
}

// NewInvitation создает новое приглашение
func (suite *SuiteWithMocks) NewInvitation(subjectID, recipientID uuid.UUID) chatt.Invitation {
	i, err := chatt.NewInvitation(subjectID, recipientID)
	suite.Require().NoError(err)
	return i
}

// AddInvitation добавляет приглашение в чат
func (suite *SuiteWithMocks) AddInvitation(chat *chatt.Chat, i chatt.Invitation) {
	suite.Require().NoError(chat.AddInvitation(i, nil))
}

// RandomString генерирует случайную строку
func RandomString2(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandomOauthToken генерирует случайный OauthToken
func (suite *SuiteWithMocks) RandomOauthToken() userr.OpenAuthToken {
	t, err := userr.NewOpenAuthToken(
		RandomString(32), // AccessToken
		"Bearer",         // TokenType
		RandomString(32), // RefreshToken
		time.Now().Add(time.Hour*24*time.Duration(rand.Intn(7)+1)), // Expiry

	)
	suite.Require().NoError(err)

	return t
}

// Генерация случайного OauthUser
func (suite *SuiteWithMocks) RandomOauthUser() userr.OpenAuthUser {
	u, err := userr.NewOpenAuthUser(
		RandomString(21),                                      // ID
		suite.Adapters.Oauth.Name(),                           // Provider
		RandomString(8)+"@example.com",                        // Email
		RandomString(6)+" "+RandomString(7),                   // Name
		"https://example.com/avatar/"+RandomString(10)+".png", // Picture
		userr.OpenAuthToken{},                                 // Token
	)
	suite.Require().NoError(err)

	return u
}

// GenerateMockUsers генерирует случайных oauth-пользователей
func (suite *SuiteWithMocks) GenerateMockUsers() map[userr.OpenAuthToken]userr.OpenAuthUser {
	tokenToUser := make(map[userr.OpenAuthToken]userr.OpenAuthUser)

	for i := 0; i < 10; i++ {
		token := suite.RandomOauthToken()
		user := suite.RandomOauthUser()
		tokenToUser[token] = user
	}

	return tokenToUser
}

// NewRndUserWithBasicAuth создает случайного пользователя с базовой аутентификацией
func (suite *SuiteWithMocks) NewRndUserWithBasicAuth() userr.User {
	user, err := userr.NewUser(gofakeit.Name(), gofakeit.Username())
	suite.Require().NoError(err)
	ba, err := userr.NewBasicAuth(gofakeit.Username()+"four", common.RndPassword())
	suite.Require().NoError(err)
	err = user.AddBasicAuth(ba)
	suite.Require().NoError(err)
	return user
}

// HasElementOfType возвращает true, если в срезе есть элемент заданного типа
func HasElementOfType2[T any](e []any) bool {
	for _, e := range e {
		if _, ok := e.(T); ok {
			return true
		}
	}
	return false
}

func (suite *SuiteWithMocks) AssertHasEventType(ee []events.Event, eventType string) {
	suite.T().Helper()
	suite.True(slices.ContainsFunc(ee, func(e events.Event) bool {
		return e.Type == eventType
	}))
}
