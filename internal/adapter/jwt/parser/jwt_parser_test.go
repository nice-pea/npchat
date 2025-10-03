package jwtParser

import (
	"testing"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jwt2 "github.com/nice-pea/npchat/internal/adapter/jwt"
	redisRegistry "github.com/nice-pea/npchat/internal/adapter/jwt/repository/redis"
)

// newTestParser - вспомогательная функция для создания тестового парсера JWT
func newTestParser(secret string) Parser {
	return Parser{Config: jwt2.Config{SecretKey: secret}}
}

// Test_JWTParser_Parse - набор тестов для проверки функции Parse парсера JWT-токенов
func Test_JWTParser_Parse(t *testing.T) {
	t.Run("Standard Parse", func(t *testing.T) {
		t.Run("валидный jwt можно разобрать и получить данные", func(t *testing.T) {
			secret := "secret"
			parser := newTestParser(secret)

			var (
				uid = uuid.New()
				sid = uuid.New()
			)

			token, err := createJWT(secret, map[string]any{
				"UserID":    uid,
				"SessionID": sid,
			})
			require.NoError(t, err)

			claims, err := parser.Parse(token)
			require.NoError(t, err)

			assert.Equal(t, uid.String(), claims.UserID)
			assert.Equal(t, sid.String(), claims.SessionID)
		})
		t.Run("jwt существующий больше exp - невалиден", func(t *testing.T) {
			secret := "secret"
			parser := newTestParser(secret)

			token, err := createJWT(secret, map[string]any{
				"UserID":    uuid.New(),
				"SessionID": uuid.New(),
				"exp":       time.Now().Add(time.Millisecond).Unix(),
			})
			require.NoError(t, err)
			time.Sleep(2 * time.Millisecond)

			claims, err := parser.Parse(token)
			assert.ErrorIs(t, err, ErrTokenExpired)
			assert.Zero(t, claims)
		})
		t.Run("jwt существующий меньше exp - валиден", func(t *testing.T) {
			secret := "secret"
			parser := newTestParser(secret)

			token, err := createJWT(secret, map[string]any{
				"UserID":    uuid.New(),
				"SessionID": uuid.New(),
				"exp":       time.Now().Add(1000 * time.Second).Unix(),
			})
			require.NoError(t, err)

			claims, err := parser.Parse(token)
			require.NoError(t, err)
			assert.NotZero(t, claims)
		})
		t.Run("jwt существующий больше 2 минут - невалиден", func(t *testing.T) {
			secret := "secret"
			parser := newTestParser(secret)

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
			parser := newTestParser(secret)

			token := `adsafs.afsfsa.gsdsddsggd`

			claims, err := parser.Parse(token)
			assert.Error(t, err)
			assert.Zero(t, claims)
		})
		t.Run("пустой jwt", func(t *testing.T) {
			secret := "secret"
			parser := newTestParser(secret)

			claims, err := parser.Parse("")
			assert.Error(t, err)
			assert.Zero(t, claims)
		})
	})

}

// createJWT - вспомогательная функция для создания JWT-токена
func createJWT(secret string, claims map[string]any) (string, error) {
	// создаем Signer
	signer, err := jwt.NewSignerHS(jwt.HS256, []byte(secret))
	if err != nil {
		return "", err
	}

	// создаем Builder
	builder := jwt.NewBuilder(signer)

	// создаем токен
	token, err := builder.Build(claims)
	if err != nil {
		return "", err
	}

	return token.String(), nil
}

// Test_parseAndValidateJWTWithInvalidation - тесты для проверки расширенной валидации токенов с Redis
func (suite *testSuite) Test_parseAndValidateJWTWithInvalidation() {
	suite.Run("если VerifyTokenWithInvalidation = false, то будет вызываться обычная проверка jwt", func() {
		suite.Parser.Config.VerifyTokenWithInvalidation = false
		var (
			uid = uuid.New()
			sid = uuid.New()
		)

		token := suite.createJWT(suite.Parser.Config.SecretKey, map[string]any{
			"UserID":    uid,
			"SessionID": sid,
		})

		claims, err := suite.Parser.Parse(token)
		suite.Require().NoError(err)

		suite.Equal(uid.String(), claims.UserID)
		suite.Equal(sid.String(), claims.SessionID)
	})
	suite.Run("если VerifyTokenWithInvalidation = true и клиентРедис не создан, то будет вызываться обычная проверка jwt", func() {
		suite.Parser.Config.VerifyTokenWithInvalidation = true
		suite.Parser.Registry = redisRegistry.Registry{}
		var (
			uid = uuid.New()
			sid = uuid.New()
		)
		token := suite.createJWT(suite.cfg.SecretKey, map[string]any{
			"UserID":    uid,
			"SessionID": sid,
		})

		claims, err := suite.Parser.Parse(token)
		suite.Require().NoError(err)

		suite.Equal(uid.String(), claims.UserID)
		suite.Equal(sid.String(), claims.SessionID)
	})
	suite.Run("если VerifyTokenWithInvalidation = true и клиентРедис создан, то будет вызываться валидация поля Iat", func() {
		var (
			uid = uuid.New()
			sid = uuid.New()
		)

		token := suite.createJWT(suite.cfg.SecretKey, map[string]any{
			"UserID":    uid,
			"SessionID": sid,
			"iat":       time.Now().Unix(),
		})

		claims, err := suite.Parser.Parse(token)
		suite.Require().NoError(err)

		suite.Equal(uid.String(), claims.UserID)
		suite.Equal(sid.String(), claims.SessionID)
	})
	suite.Run("в кэше есть запись, анулирующее все токены данной сессии", func() {
		var (
			uid = uuid.New()
			sid = uuid.New()
		)
		token := suite.createJWT(suite.cfg.SecretKey, map[string]any{
			"UserID":    uid,
			"SessionID": sid,
			"iat":       time.Now().Unix(),
		})

		err := suite.Parser.Registry.RegisterIssueTime(sid, time.Now())
		suite.Require().NoError(err)

		claims, err := suite.Parser.Parse(token)
		suite.Require().ErrorIs(err, ErrTokenRevoked)
		suite.Assert().Zero(claims)
	})
	suite.Run("можно заранее анулировать будущие токены", func() {
		var (
			uid = uuid.New()
			sid = uuid.New()
		)
		err := suite.Parser.Registry.RegisterIssueTime(sid, time.Now().Add(time.Hour))
		suite.Require().NoError(err)

		token := suite.createJWT(suite.cfg.SecretKey, map[string]any{
			"UserID":    uid,
			"SessionID": sid,
			"iat":       time.Now().Unix(),
		})

		claims, err := suite.Parser.Parse(token)
		suite.Require().ErrorIs(err, ErrTokenRevoked)
		suite.Assert().Zero(claims)
	})
	suite.Run("если токен создан после даты анулирования, то токен действителен", func() {
		var (
			uid = uuid.New()
			sid = uuid.New()
		)

		err := suite.Parser.Registry.RegisterIssueTime(sid, time.Now().Add(-time.Hour))
		suite.Require().NoError(err)

		token := suite.createJWT(suite.cfg.SecretKey, map[string]any{
			"UserID":    uid,
			"SessionID": sid,
			"iat":       time.Now().Unix(),
		})

		claims, err := suite.Parser.Parse(token)
		suite.Require().NoError(err)

		suite.Equal(uid.String(), claims.UserID)
		suite.Equal(sid.String(), claims.SessionID)
	})
	suite.Run("если iat пустой то вернется ошибка ErrIatEmpty", func() {
		var (
			uid = uuid.New()
			sid = uuid.New()
		)

		token := suite.createJWT(suite.cfg.SecretKey, map[string]any{
			"UserID":    uid,
			"SessionID": sid,
		})

		claims, err := suite.Parser.Parse(token)
		suite.Require().ErrorIs(err, ErrIatEmpty)
		suite.Zero(claims)
	})
	suite.Run("анулируются только те токены, у которых iat меньше даты в кэше", func() {
		var (
			uid = uuid.New()
			sid = uuid.New()
		)
		secret := suite.cfg.SecretKey
		token1 := suite.createJWT(secret, map[string]any{
			"UserID":    uid,
			"SessionID": sid,
			"iat":       time.Now().Unix(),
		})

		err := suite.Parser.Registry.RegisterIssueTime(sid, time.Now())
		suite.Require().NoError(err)

		time.Sleep(time.Second)
		token2 := suite.createJWT(secret, map[string]any{
			"UserID":    uid,
			"SessionID": sid,
			"iat":       time.Now().Unix(),
		})

		claims, err := suite.Parser.Parse(token1)
		suite.Require().ErrorIs(err, ErrTokenRevoked)
		suite.Assert().Zero(claims)

		claims, err = suite.Parser.Parse(token2)
		suite.Require().NoError(err)
		suite.Equal(sid.String(), claims.SessionID)
		suite.Equal(uid.String(), claims.UserID)
	})
}
