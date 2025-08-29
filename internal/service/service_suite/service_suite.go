package servicesuite

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
	"github.com/nice-pea/npchat/internal/service"
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
		Oauth service.OAuthProvider
	}
	mockOAuthTokens map[string]userr.OpenAuthToken
	mockOAuthUsers  map[userr.OpenAuthToken]userr.OpenAuthUser
}

var pgsqlDSN = os.Getenv("TEST_PGSQL_DSN")

func Test_ServicesTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testifySuite.Run(t, new(Suite))
}

func (suite *Suite) newPgsqlExternalFactory(dsn string) (*pgsqlRepository.Factory, func()) {
	factory, err := pgsqlRepository.InitFactory(pgsqlRepository.Config{
		DSN: dsn,
	})
	suite.Require().NoError(err)

	return factory, func() { _ = factory.Close() }
}

func (suite *Suite) newPgsqlContainerFactory() (f *pgsqlRepository.Factory, closer func()) {
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
	suite.mockOAuthUsers = suite.generateMockUsers()
	suite.mockOAuthTokens = make(map[string]userr.OpenAuthToken, len(suite.mockOAuthUsers))
	for token := range suite.mockOAuthUsers {
		suite.mockOAuthTokens[randomString(13)] = token
	}
	suite.Adapters.Oauth = &oauthProvider.Mock{
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

// upsertChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *Suite) upsertChat(chat chatt.Chat) chatt.Chat {
	err := suite.RR.Chats.Upsert(chat)
	suite.Require().NoError(err)

	return chat
}

// upsertChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *Suite) rndChat() chatt.Chat {
	chat, err := chatt.NewChat(gofakeit.Noun(), uuid.New())
	suite.Require().NoError(err)

	return chat
}

// newParticipant создает случайного участника
func (suite *Suite) newParticipant(userID uuid.UUID) chatt.Participant {
	p, err := chatt.NewParticipant(userID)
	suite.Require().NoError(err)
	return p
}

func (suite *Suite) addRndParticipant(chat *chatt.Chat) chatt.Participant {
	p, err := chatt.NewParticipant(uuid.New())
	suite.Require().NoError(err)
	suite.Require().NoError(chat.AddParticipant(p))

	return p
}

func (suite *Suite) addParticipant(chat *chatt.Chat, p chatt.Participant) {
	suite.Require().NoError(chat.AddParticipant(p))
}

func (suite *Suite) newInvitation(subjectID, recipientID uuid.UUID) chatt.Invitation {
	i, err := chatt.NewInvitation(subjectID, recipientID)
	suite.Require().NoError(err)
	return i
}

func (suite *Suite) addInvitation(chat *chatt.Chat, i chatt.Invitation) {
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
func (suite *Suite) randomOAuthToken() userr.OpenAuthToken {
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
func (suite *Suite) randomOAuthUser() userr.OpenAuthUser {
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

func (suite *Suite) generateMockUsers() map[userr.OpenAuthToken]userr.OpenAuthUser {
	tokenToUser := make(map[userr.OpenAuthToken]userr.OpenAuthUser)

	for i := 0; i < 10; i++ {
		token := suite.randomOAuthToken()
		user := suite.randomOAuthUser()
		tokenToUser[token] = user
	}

	return tokenToUser
}

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
