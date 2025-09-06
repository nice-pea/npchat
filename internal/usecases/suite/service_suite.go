package serviceSuite

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	oauthProvider "github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/common"
	"github.com/nice-pea/npchat/internal/domain/chatt"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	pgsqlRepository "github.com/nice-pea/npchat/internal/repository/pgsql_repository"
	"github.com/nice-pea/npchat/internal/usecases/events"
	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
)

type Suite struct {
	testifySuite.Suite
	factory       *pgsqlRepository.Factory
	factoryCloser func()
	RR            struct {
		Chats    chatt.Repository
		Sessions sessionn.Repository
		Users    userr.Repository
	}
	Adapters struct {
		Oauth oauth.OAuthProvider
	}
	MockOAuthTokens map[string]userr.OpenAuthToken
	MockOAuthUsers  map[userr.OpenAuthToken]userr.OpenAuthUser
}

var (
	pgsqlDSN = os.Getenv("TEST_PGSQL_DSN")
)

// newPgsqlExternalFactory создает фабрику репозиториев  для тестирования, реализованных с помощью подключения к postgres по DSN
func (suite *Suite) newPgsqlExternalFactory(dsn string) (*pgsqlRepository.Factory, func()) {
	factory, err := pgsqlRepository.InitFactory(pgsqlRepository.Config{
		DSN: dsn,
	})
	suite.Require().NoError(err)

	return factory, func() { _ = factory.Close() }
}

// newPgsqlContainerFactory создает фабрику репозиториев  для тестирования, реализованных с помощью postgres контейнеров
func (suite *Suite) newPgsqlContainerFactory() (f *pgsqlRepository.Factory, closer func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Поиск скрптов с миграциями
	_, b, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(b), "../../../infra/pgsql/init/*.up.sql")
	migrations, err := filepath.Glob(migrationsDir)
	suite.Require().NoError(err)
	suite.Require().NotZero(migrations)

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17",
		postgres.WithInitScripts(migrations...),
		postgres.WithDatabase("test_npc_db"),
		postgres.WithUsername("test_npc_user"),
		postgres.WithPassword("test_npc_password"),
		postgres.BasicWaitStrategies(),
	)
	suite.Require().NoError(err)
	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	suite.Require().NoError(err)

	f, err = pgsqlRepository.InitFactory(pgsqlRepository.Config{
		DSN: dsn,
	})
	suite.Require().NoError(err)

	return f, func() {
		_ = f.Close()
		_ = postgresContainer.Terminate(context.Background())
	}
}

// SetupTest выполняется перед каждым тестом, связанным с suite
func (suite *Suite) SetupTest() {
	// Инициализация фабрики репозиториев
	if pgsqlDSN != "" {
		suite.factory, suite.factoryCloser = suite.newPgsqlExternalFactory(pgsqlDSN)
	} else {
		suite.factory, suite.factoryCloser = suite.newPgsqlContainerFactory()
	}

	// Инициализация репозиториев
	suite.RR.Chats = suite.factory.NewChattRepository()
	suite.RR.Users = suite.factory.NewUserrRepository()
	suite.RR.Sessions = suite.factory.NewSessionnRepository()

	// Инициализация адаптеров
	suite.MockOAuthUsers = suite.GenerateMockUsers()
	suite.MockOAuthTokens = make(map[string]userr.OpenAuthToken, len(suite.MockOAuthUsers))
	for token := range suite.MockOAuthUsers {
		suite.MockOAuthTokens[RandomString(13)] = token
	}
	suite.Adapters.Oauth = &oauthProvider.Mock{
		ExchangeFunc: func(code string) (userr.OpenAuthToken, error) {
			token, ok := suite.MockOAuthTokens[code]
			if !ok {
				return userr.OpenAuthToken{}, errors.New("token not found")
			}
			return token, nil
		},
		UserFunc: func(token userr.OpenAuthToken) (userr.OpenAuthUser, error) {
			user, ok := suite.MockOAuthUsers[token]
			if !ok {
				return userr.OpenAuthUser{}, errors.New("user not found")
			}
			return user, nil
		},
		AuthorizationURLFunc: func(state string) string {
			return "https://provider.com/o/oauth2/auth?code=somecode&state=" + state
		},
	}
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *Suite) TearDownSubTest() {
	err := suite.factory.Cleanup()
	suite.Require().NoError(err)
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *Suite) TearDownTest() {
	suite.factoryCloser()
}

// EqualSessions сравнивает две сессии
func (suite *Suite) EqualSessions(s1, s2 sessionn.Session) {
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
func (suite *Suite) UpsertChat(chat chatt.Chat) chatt.Chat {
	err := suite.RR.Chats.Upsert(chat)
	suite.Require().NoError(err)

	return chat
}

// RndChat создает случайный чат
func (suite *Suite) RndChat() chatt.Chat {
	chat, err := chatt.NewChat(gofakeit.Noun(), uuid.New(), nil)
	suite.Require().NoError(err)

	return chat
}

// NewParticipant создает случайного участника
func (suite *Suite) NewParticipant(userID uuid.UUID) chatt.Participant {
	p, err := chatt.NewParticipant(userID)
	suite.Require().NoError(err)
	return p
}

// AddRndParticipant добавляет случайного участника в чат
func (suite *Suite) AddRndParticipant(chat *chatt.Chat) chatt.Participant {
	p, err := chatt.NewParticipant(uuid.New())
	suite.Require().NoError(err)
	suite.Require().NoError(chat.AddParticipant(p, nil))

	return p
}

// AddParticipant добавляет участника в чат
func (suite *Suite) AddParticipant(chat *chatt.Chat, p chatt.Participant) {
	suite.Require().NoError(chat.AddParticipant(p, nil))
}

// NewInvitation создает новое приглашение
func (suite *Suite) NewInvitation(subjectID, recipientID uuid.UUID) chatt.Invitation {
	i, err := chatt.NewInvitation(subjectID, recipientID)
	suite.Require().NoError(err)
	return i
}

// AddInvitation добавляет приглашение в чат
func (suite *Suite) AddInvitation(chat *chatt.Chat, i chatt.Invitation) {
	suite.Require().NoError(chat.AddInvitation(i, nil))
}

// RandomString генерирует случайную строку
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandomOAuthToken генерирует случайный OAuthToken
func (suite *Suite) RandomOAuthToken() userr.OpenAuthToken {
	t, err := userr.NewOpenAuthToken(
		RandomString(32), // AccessToken
		"Bearer",         // TokenType
		RandomString(32), // RefreshToken
		time.Now().Add(time.Hour*24*time.Duration(rand.Intn(7)+1)), // Expiry

	)
	suite.Require().NoError(err)

	return t
}

// Генерация случайного OAuthUser
func (suite *Suite) RandomOAuthUser() userr.OpenAuthUser {
	u, err := userr.NewOpenAuthUser(
		RandomString(21),                                      // ID
		(&oauthProvider.Mock{}).Name(),                        // Provider
		RandomString(8)+"@example.com",                        // Email
		RandomString(6)+" "+RandomString(7),                   // Name
		"https://example.com/avatar/"+RandomString(10)+".png", // Picture
		userr.OpenAuthToken{},                                 // Token
	)
	suite.Require().NoError(err)

	return u
}

// GenerateMockUsers генерирует случайных oauth-пользователей
func (suite *Suite) GenerateMockUsers() map[userr.OpenAuthToken]userr.OpenAuthUser {
	tokenToUser := make(map[userr.OpenAuthToken]userr.OpenAuthUser)

	for i := 0; i < 10; i++ {
		token := suite.RandomOAuthToken()
		user := suite.RandomOAuthUser()
		tokenToUser[token] = user
	}

	return tokenToUser
}

// NewRndUserWithBasicAuth создает случайного пользователя с базовой аутентификацией
func (suite *Suite) NewRndUserWithBasicAuth() userr.User {
	user, err := userr.NewUser(gofakeit.Name(), gofakeit.Username())
	suite.Require().NoError(err)
	ba, err := userr.NewBasicAuth(gofakeit.Username()+"four", common.RndPassword())
	suite.Require().NoError(err)
	err = user.AddBasicAuth(ba)
	suite.Require().NoError(err)
	return user
}

// HasElementOfType возвращает true, если в срезе есть элемент заданного типа
func HasElementOfType[T any](e []any) bool {
	for _, e := range e {
		if _, ok := e.(T); ok {
			return true
		}
	}
	return false
}

func (suite *Suite) AssertHasEventType(ee []events.Event, eventType string) {
	suite.T().Helper()
	suite.True(slices.ContainsFunc(ee, func(e events.Event) bool {
		return e.Type == eventType
	}))
}
