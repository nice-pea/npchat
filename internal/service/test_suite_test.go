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
func (suite *testSuite) TearDownSubTest() {
	err := suite.factory.Cleanup()
	suite.Require().NoError(err)
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *testSuite) TearDownTest() {
	suite.factoryCloser()
}

// upsertChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *testSuite) upsertChat(chat chatt.Chat) chatt.Chat {
	err := suite.rr.chats.Upsert(chat)
	suite.Require().NoError(err)

	return chat
}

// upsertChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *testSuite) rndChat() chatt.Chat {
	chat, err := chatt.NewChat(gofakeit.Noun(), uuid.New())
	suite.Require().NoError(err)

	return chat
}

// newParticipant создает случайного участника
func (suite *testSuite) newParticipant(userID uuid.UUID) chatt.Participant {
	p, err := chatt.NewParticipant(userID)
	suite.Require().NoError(err)
	return p
}

func (suite *testSuite) addRndParticipant(chat *chatt.Chat) chatt.Participant {
	p, err := chatt.NewParticipant(uuid.New())
	suite.Require().NoError(err)
	suite.Require().NoError(chat.AddParticipant(p))

	return p
}

func (suite *testSuite) addParticipant(chat *chatt.Chat, p chatt.Participant) {
	suite.Require().NoError(chat.AddParticipant(p))
}

func (suite *testSuite) newInvitation(subjectID, recipientID uuid.UUID) chatt.Invitation {
	i, err := chatt.NewInvitation(subjectID, recipientID)
	suite.Require().NoError(err)
	return i
}

func (suite *testSuite) addInvitation(chat *chatt.Chat, i chatt.Invitation) {
	suite.Require().NoError(chat.AddInvitation(i))
}

// randomString генерирует случайную строку
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// randomOAuthToken генерирует случайный OAuthToken
func (suite *testSuite) randomOAuthToken() userr.OpenAuthToken {
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
func (suite *testSuite) randomOAuthUser() userr.OpenAuthUser {
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

func (suite *testSuite) generateMockUsers() map[userr.OpenAuthToken]userr.OpenAuthUser {
	tokenToUser := make(map[userr.OpenAuthToken]userr.OpenAuthUser)

	for i := 0; i < 10; i++ {
		token := suite.randomOAuthToken()
		user := suite.randomOAuthUser()
		tokenToUser[token] = user
	}

	return tokenToUser
}
