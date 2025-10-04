package jwtParser

import (
	"testing"
	"time"

	"github.com/nice-pea/npchat/internal/adapter/jwt"
	mockJwtParser "github.com/nice-pea/npchat/internal/adapter/jwt/parser/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jwt2 "github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"
)

// Test_JWTParser_Parse набор тестов для проверки функции Parse парсера JWT-токенов
func Test_JWTParser_Parse(t *testing.T) {
	t.Run("Парсер без Registry", func(t *testing.T) {
		t.Run("валидный jwt можно разобрать и получить данные", func(t *testing.T) {
			secret := "secret"
			parser := parserWithoutRegistry(t, secret)

			var (
				uid = uuid.New()
				sid = uuid.New()
			)

			token := createJWT(t, secret, map[string]any{
				"UserID":    uid,
				"SessionID": sid,
			})

			claims, err := parser.Parse(token)
			require.NoError(t, err)

			assert.Equal(t, uid.String(), claims.UserID)
			assert.Equal(t, sid.String(), claims.SessionID)
		})

		t.Run("jwt существующий больше exp - невалиден", func(t *testing.T) {
			secret := "secret"
			parser := parserWithoutRegistry(t, secret)

			token := createJWT(t, secret, map[string]any{
				"UserID":    uuid.New(),
				"SessionID": uuid.New(),
				"exp":       time.Now().Add(time.Millisecond).Unix(),
			})
			time.Sleep(2 * time.Millisecond)

			claims, err := parser.Parse(token)
			assert.ErrorIs(t, err, ErrTokenExpired)
			assert.Zero(t, claims)
		})

		t.Run("jwt существующий меньше exp - валиден", func(t *testing.T) {
			secret := "secret"
			parser := parserWithoutRegistry(t, secret)

			token := createJWT(t, secret, map[string]any{
				"UserID":    uuid.New(),
				"SessionID": uuid.New(),
				"exp":       time.Now().Add(1000 * time.Second).Unix(),
			})

			claims, err := parser.Parse(token)
			require.NoError(t, err)
			assert.NotZero(t, claims)
		})

		t.Run("jwt существующий больше 2 минут - невалиден", func(t *testing.T) {
			secret := "secret"
			parser := parserWithoutRegistry(t, secret)

			//token - истекший jwt токен
			// содержит данные:
			/* {
			  "SessionID": "456",
			  "UserID": "123",
			  "exp": 1759395766
			} */
			token := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJTZXNzaW9uSUQiOiI0NTYiLCJVc2VySUQiOiIxMjMiLCJleHAiOjE3NTkzOTU3NjZ9.SYQurl5gsOt42K2d0Vyp-RuZluANRNuGMyUNd6RfWtk`

			claims, err := parser.Parse(token)
			assert.ErrorIs(t, err, ErrTokenExpired)
			assert.Zero(t, claims)
		})

		t.Run("невалидный jwt", func(t *testing.T) {
			secret := "secret"
			parser := parserWithoutRegistry(t, secret)

			token := `adsafs.afsfsa.gsdsddsggd`
			claims, err := parser.Parse(token)
			require.Error(t, err)
			assert.Zero(t, claims)
		})

		t.Run("пустой jwt", func(t *testing.T) {
			secret := "secret"
			parser := parserWithoutRegistry(t, secret)

			claims, err := parser.Parse("")
			assert.Error(t, err)
			assert.Zero(t, claims)
		})
	})
	t.Run("Парсер c Registry", func(t *testing.T) {
		t.Run("если VerifyTokenWithInvalidation = false, то будет вызываться обычная проверка jwt", func(t *testing.T) {
			parser, _ := parserWithRegistry(t)
			parser.Config.VerifyTokenWithInvalidation = false
			var (
				uid = uuid.New()
				sid = uuid.New()
			)

			token := createJWT(t, parser.Config.SecretKey, map[string]any{
				"UserID":    uid,
				"SessionID": sid,
			})

			claims, err := parser.Parse(token)
			require.NoError(t, err)

			assert.Equal(t, uid.String(), claims.UserID)
			assert.Equal(t, sid.String(), claims.SessionID)
		})

		t.Run("если VerifyTokenWithInvalidation = true и клиентРедис не создан, то будет вызываться обычная проверка jwt", func(t *testing.T) {
			parser, _ := parserWithRegistry(t)
			parser.Config.VerifyTokenWithInvalidation = true
			parser.Registry = nil
			var (
				uid = uuid.New()
				sid = uuid.New()
			)
			token := createJWT(t, parser.Config.SecretKey, map[string]any{
				"UserID":    uid,
				"SessionID": sid,
			})

			claims, err := parser.Parse(token)
			require.NoError(t, err)

			assert.Equal(t, uid.String(), claims.UserID)
			assert.Equal(t, sid.String(), claims.SessionID)
		})

		t.Run("если VerifyTokenWithInvalidation = true и клиентРедис создан, то будет вызываться валидация поля Iat", func(t *testing.T) {
			parser, registryMock := parserWithRegistry(t)
			var (
				uid = uuid.New()
				sid = uuid.New()
			)

			token := createJWT(t, parser.Config.SecretKey, map[string]any{
				"UserID":    uid,
				"SessionID": sid,
				"iat":       time.Now().Unix(),
			})

			registryMock.EXPECT().IssueTime(sid).Return(time.Time{}, nil)

			claims, err := parser.Parse(token)
			require.NoError(t, err)

			assert.Equal(t, uid.String(), claims.UserID)
			assert.Equal(t, sid.String(), claims.SessionID)
		})

		t.Run("если токен создан после даты анулирования, то токен действителен", func(t *testing.T) {
			parser, registryMock := parserWithRegistry(t)
			var (
				uid = uuid.New()
				sid = uuid.New()
			)

			registryMock.EXPECT().IssueTime(sid).Return(time.Now().Add(-time.Hour), nil)

			token := createJWT(t, parser.Config.SecretKey, map[string]any{
				"UserID":    uid,
				"SessionID": sid,
				"iat":       time.Now().Unix(),
			})

			claims, err := parser.Parse(token)
			require.NoError(t, err)

			assert.Equal(t, uid.String(), claims.UserID)
			assert.Equal(t, sid.String(), claims.SessionID)
		})

		t.Run("если iat пустой то вернется ошибка ErrIatEmpty", func(t *testing.T) {
			parser, registryMock := parserWithRegistry(t)
			var (
				uid = uuid.New()
				sid = uuid.New()
			)

			token := createJWT(t, parser.Config.SecretKey, map[string]any{
				"UserID":    uid,
				"SessionID": sid,
			})
			registryMock.EXPECT().IssueTime(sid).Return(time.Now(), nil)

			claims, err := parser.Parse(token)
			require.ErrorIs(t, err, ErrIatEmpty)
			assert.Zero(t, claims)
		})

		t.Run("анулируются только те токены, у которых iat меньше даты в кэше", func(t *testing.T) {
			parser, registryMock := parserWithRegistry(t)
			var (
				uid = uuid.New()
				sid = uuid.New()
			)
			secret := parser.Config.SecretKey
			token1 := createJWT(t, secret, map[string]any{
				"UserID":    uid,
				"SessionID": sid,
				"iat":       time.Now().Unix(),
			})

			timeNow := time.Now()

			time.Sleep(time.Second)
			token2 := createJWT(t, secret, map[string]any{
				"UserID":    uid,
				"SessionID": sid,
				"iat":       time.Now().Unix(),
			})

			registryMock.EXPECT().IssueTime(sid).Return(timeNow, nil)
			claims, err := parser.Parse(token1)
			require.ErrorIs(t, err, ErrTokenRevoked)
			assert.Zero(t, claims)

			claims, err = parser.Parse(token2)
			require.NoError(t, err)
			assert.Equal(t, sid.String(), claims.SessionID)
			assert.Equal(t, uid.String(), claims.UserID)
		})
	})
}

// parserWithRegistry создает парсер с Registry
func parserWithRegistry(t *testing.T) (Parser, *mockJwtParser.Registry) {
	// создаем mockRegistry
	registryMock := mockJwtParser.NewRegistry(t)

	// создаем Parser
	cfg := jwt.Config{
		SecretKey:                   "secret",
		VerifyTokenWithInvalidation: true,
	}
	return Parser{
		Config:   cfg,
		Registry: registryMock,
	}, registryMock
}

// parserWithoutRegistry создает парсер без Registry
func parserWithoutRegistry(t *testing.T, secret string) Parser {
	return Parser{Config: jwt.Config{SecretKey: secret}}
}

// createJWT создает JWT токен для тестов
func createJWT(t *testing.T, secret string, claims map[string]any) string {
	// создаем Signer
	signer, err := jwt2.NewSignerHS(jwt2.HS256, []byte(secret))
	require.NoError(t, err)
	// создаем Builder
	builder := jwt2.NewBuilder(signer)

	// создаем токен
	token, err := builder.Build(claims)
	require.NoError(t, err)

	return token.String()
}
