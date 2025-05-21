package service

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/saime-0/nice-pea-chat/internal/adapter"
	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
)

type servicesTestSuite struct {
	suite.Suite
	factory *sqlite.RepositoryFactory
	rr      struct {
		chats         domain.ChatsRepository
		members       domain.MembersRepository
		invitations   domain.InvitationsRepository
		sessions      domain.SessionsRepository
		users         domain.UsersRepository
		authnPassword domain.AuthnPasswordRepository
		oauth         domain.OAuthRepository
	}
	ss struct {
		chats         *Chats
		members       *Members
		invitations   *Invitations
		sessions      *Sessions
		authnPassword *AuthnPassword
		oauth         *OAuth
	}
	ad struct {
		oauth adapter.OAuthGoogle
	}
	mockOauthCodes  map[string]domain.OAuthToken
	mockGoogleUsers map[domain.OAuthToken]domain.OAuthGoogleUser
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
	suite.rr.members = suite.factory.NewMembersRepository()
	suite.rr.invitations = suite.factory.NewInvitationsRepository()
	suite.rr.users = suite.factory.NewUsersRepository()
	suite.rr.sessions = suite.factory.NewSessionsRepository()
	suite.rr.authnPassword = suite.factory.NewAuthnPasswordRepository()
	suite.rr.oauth = suite.factory.NewOAuthRepository()

	// Инициализация адаптеров
	suite.mockGoogleUsers = generateMockUsers()
	suite.mockOauthCodes = make(map[string]domain.OAuthToken, len(suite.mockGoogleUsers))
	for token := range suite.mockGoogleUsers {
		suite.mockOauthCodes[randomString(13)] = token
	}
	suite.ad.oauth = &adapter.MockOAuthGoogle{
		ExchangeFunc: func(code string) (domain.OAuthToken, error) {
			token, ok := suite.mockOauthCodes[code]
			if !ok {
				return domain.OAuthToken{}, errors.New("token not found")
			}
			return token, nil
		},
		UserFunc: func(token domain.OAuthToken) (domain.OAuthGoogleUser, error) {
			user, ok := suite.mockGoogleUsers[token]
			if !ok {
				return domain.OAuthGoogleUser{}, errors.New("user not found")
			}
			return user, nil
		},
		AuthCodeURLFunc: func(state string) string {
			return "https://accounts.google.com/o/oauth2/auth?state=" + state
		},
	}

	// Создание сервисов
	suite.ss.chats = &Chats{
		ChatsRepo:   suite.rr.chats,
		MembersRepo: suite.rr.members,
	}
	suite.ss.members = &Members{
		ChatsRepo:   suite.rr.chats,
		MembersRepo: suite.rr.members,
	}
	suite.ss.invitations = &Invitations{
		ChatsRepo:       suite.rr.chats,
		MembersRepo:     suite.rr.members,
		InvitationsRepo: suite.rr.invitations,
		UsersRepo:       suite.rr.users,
	}
	suite.ss.sessions = &Sessions{
		SessionsRepo: suite.rr.sessions,
	}
	suite.ss.authnPassword = &AuthnPassword{
		AuthnPasswordRepo: suite.rr.authnPassword,
		SessionsRepo:      suite.rr.sessions,
		UsersRepo:         suite.rr.users,
	}
	suite.ss.oauth = &OAuth{
		Google:    suite.ad.oauth,
		OAuthRepo: suite.rr.oauth,
	}
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *servicesTestSuite) TearDownSubTest() {
	err := suite.factory.Close()
	suite.Require().NoError(err)
}

// saveChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveChat(chat domain.Chat) domain.Chat {
	err := suite.rr.chats.Save(chat)
	suite.Require().NoError(err)

	return chat
}

// saveMember сохраняет участника в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveMember(member domain.Member) domain.Member {
	err := suite.rr.members.Save(member)
	suite.Require().NoError(err)

	return member
}

// saveInvitation сохраняет приглашение в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveInvitation(invitation domain.Invitation) domain.Invitation {
	err := suite.rr.invitations.Save(invitation)
	suite.Require().NoError(err)

	return invitation
}

// saveUser сохраняет пользователя в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveUser(user domain.User) domain.User {
	err := suite.rr.users.Save(user)
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
func randomOAuthToken() domain.OAuthToken {
	return domain.OAuthToken{
		AccessToken:  randomString(32),
		TokenType:    "Bearer",
		RefreshToken: randomString(32),
		Expiry:       time.Now().Add(time.Hour * 24 * time.Duration(rand.Intn(7)+1)),
	}
}

// Генерация случайного OAuthGoogleUser
func randomOAuthGoogleUser() domain.OAuthGoogleUser {
	return domain.OAuthGoogleUser{
		ID:            randomString(21), // Google ID обычно длина ~21
		Email:         randomString(8) + "@example.com",
		VerifiedEmail: true,
		Name:          randomString(6) + " " + randomString(7),
		GivenName:     randomString(6),
		FamilyName:    randomString(7),
		Picture:       "https://example.com/avatar/" + randomString(10) + ".png",
		Locale:        "en",
	}
}

// Инициализация карты tokenToGoogleUser
func generateMockUsers() map[domain.OAuthToken]domain.OAuthGoogleUser {
	tokenToGoogleUser := make(map[domain.OAuthToken]domain.OAuthGoogleUser)

	for i := 0; i < 10; i++ {
		token := randomOAuthToken()
		user := randomOAuthGoogleUser()
		tokenToGoogleUser[token] = user
	}

	return tokenToGoogleUser
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
