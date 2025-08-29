package service

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	oauthProvider "github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/domain/chatt"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	pgsqlRepository "github.com/nice-pea/npchat/internal/repository/pgsql_repository"
	sessionsFind "github.com/nice-pea/npchat/internal/service/sessions/find"
)

type testSuite struct {
	testifySuite.Suite
	factory       *pgsqlRepository.Factory
	factoryCloser func()
	rr            struct {
		chats    chatt.Repository
		sessions sessionn.Repository
		users    userr.Repository
	}
	ss struct {
		chats        *Chats
		sessionsFind *sessionsFind.SessionsFindUsecase
		users        *Users
	}
	ad struct {
		oauth OAuthProvider
	}
	mockOAuthTokens map[string]userr.OpenAuthToken
	mockOAuthUsers  map[userr.OpenAuthToken]userr.OpenAuthUser
}

var pgsqlDSN = os.Getenv("TEST_PGSQL_DSN")

func Test_ServicesTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) newPgsqlExternalFactory(dsn string) (*pgsqlRepository.Factory, func()) {
	factory, err := pgsqlRepository.InitFactory(pgsqlRepository.Config{
		DSN: dsn,
	})
	suite.Require().NoError(err)

	return factory, func() { _ = factory.Close() }
}

func (suite *testSuite) newPgsqlContainerFactory() (f *pgsqlRepository.Factory, closer func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find all migrations
	migrations, err := filepath.Glob("../../infra/pgsql/init/*.sql")
	suite.Require().NoError(err)

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
func (suite *testSuite) SetupTest() {
	// Инициализация фабрики репозиториев
	if pgsqlDSN != "" {
		suite.factory, suite.factoryCloser = suite.newPgsqlExternalFactory(pgsqlDSN)
	} else {
		suite.factory, suite.factoryCloser = suite.newPgsqlContainerFactory()
	}

	// Инициализация репозиториев
	suite.rr.chats = suite.factory.NewChattRepository()
	suite.rr.users = suite.factory.NewUserrRepository()
	suite.rr.sessions = suite.factory.NewSessionnRepository()

	// Инициализация адаптеров
	suite.mockOAuthUsers = suite.GenerateMockUsers()
	suite.mockOAuthTokens = make(map[string]userr.OpenAuthToken, len(suite.mockOAuthUsers))
	for token := range suite.mockOAuthUsers {
		suite.mockOAuthTokens[RandomString(13)] = token
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
	suite.ss.sessionsFind = &sessionsFind.SessionsFindUsecase{
		Repo: suite.rr.sessions,
	}
	suite.ss.users = &Users{
		Providers:    OAuthProviders{suite.ad.oauth.Name(): suite.ad.oauth},
		Repo:         suite.rr.users,
		SessionsRepo: suite.rr.sessions,
	}
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *testSuite) TearDownSubTest() {
	err := suite.factory.Cleanup()
	suite.Require().NoError(err)
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *testSuite) TearDownTest() {
	suite.factoryCloser()
}

// UpsertChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *testSuite) UpsertChat(chat chatt.Chat) chatt.Chat {
	err := suite.rr.chats.Upsert(chat)
	suite.Require().NoError(err)

	return chat
}

// upsertChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *testSuite) RndChat() chatt.Chat {
	chat, err := chatt.NewChat(gofakeit.Noun(), uuid.New())
	suite.Require().NoError(err)

	return chat
}

// NewParticipant создает случайного участника
func (suite *testSuite) NewParticipant(userID uuid.UUID) chatt.Participant {
	p, err := chatt.NewParticipant(userID)
	suite.Require().NoError(err)
	return p
}

func (suite *testSuite) AddRndParticipant(chat *chatt.Chat) chatt.Participant {
	p, err := chatt.NewParticipant(uuid.New())
	suite.Require().NoError(err)
	suite.Require().NoError(chat.AddParticipant(p))

	return p
}

func (suite *testSuite) AddParticipant(chat *chatt.Chat, p chatt.Participant) {
	suite.Require().NoError(chat.AddParticipant(p))
}

func (suite *testSuite) NewInvitation(subjectID, recipientID uuid.UUID) chatt.Invitation {
	i, err := chatt.NewInvitation(subjectID, recipientID)
	suite.Require().NoError(err)
	return i
}

func (suite *testSuite) AddInvitation(chat *chatt.Chat, i chatt.Invitation) {
	suite.Require().NoError(chat.AddInvitation(i))
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
func (suite *testSuite) RandomOAuthToken() userr.OpenAuthToken {
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
func (suite *testSuite) RandomOAuthUser() userr.OpenAuthUser {
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

func (suite *testSuite) GenerateMockUsers() map[userr.OpenAuthToken]userr.OpenAuthUser {
	tokenToUser := make(map[userr.OpenAuthToken]userr.OpenAuthUser)

	for i := 0; i < 10; i++ {
		token := suite.RandomOAuthToken()
		user := suite.RandomOAuthUser()
		tokenToUser[token] = user
	}

	return tokenToUser
}
