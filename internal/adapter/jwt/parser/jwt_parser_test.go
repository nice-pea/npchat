package jwtParser

import (
	"time"

	"github.com/google/uuid"
)

// Test_JWTParser_Parse набор тестов для проверки функции Parse парсера JWT-токенов
func (suite *testSuite) Test_JWTParser_Parse() {
	suite.Run("Парсер без Registry", func() {
		suite.Run("валидный jwt можно разобрать и получить данные", func() {
			secret := "secret"
			parser := suite.parserWithoutRegistry(secret)

			var (
				uid = uuid.New()
				sid = uuid.New()
			)

			token := suite.createJWT(secret, map[string]any{
				"UserID":    uid,
				"SessionID": sid,
			})

			claims, err := parser.Parse(token)
			suite.Require().NoError(err)

			suite.Equal(uid.String(), claims.UserID)
			suite.Equal(sid.String(), claims.SessionID)
		})

		suite.Run("jwt существующий больше exp - невалиден", func() {
			secret := "secret"
			parser := suite.parserWithoutRegistry(secret)

			token := suite.createJWT(secret, map[string]any{
				"UserID":    uuid.New(),
				"SessionID": uuid.New(),
				"exp":       time.Now().Add(time.Millisecond).Unix(),
			})
			time.Sleep(2 * time.Millisecond)

			claims, err := parser.Parse(token)
			suite.ErrorIs(err, ErrTokenExpired)
			suite.Zero(claims)
		})

		suite.Run("jwt существующий меньше exp - валиден", func() {
			secret := "secret"
			parser := suite.parserWithoutRegistry(secret)

			token := suite.createJWT(secret, map[string]any{
				"UserID":    uuid.New(),
				"SessionID": uuid.New(),
				"exp":       time.Now().Add(1000 * time.Second).Unix(),
			})

			claims, err := parser.Parse(token)
			suite.Require().NoError(err)
			suite.NotZero(claims)
		})

		suite.Run("jwt существующий больше 2 минут - невалиден", func() {
			secret := "secret"
			parser := suite.parserWithoutRegistry(secret)

			//token - истекший jwt токен
			// содержит данные:
			/* {
			  "SessionID": "456",
			  "UserID": "123",
			  "exp": 1759395766
			} */
			token := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJTZXNzaW9uSUQiOiI0NTYiLCJVc2VySUQiOiIxMjMiLCJleHAiOjE3NTkzOTU3NjZ9.SYQurl5gsOt42K2d0Vyp-RuZluANRNuGMyUNd6RfWtk`

			claims, err := parser.Parse(token)
			suite.ErrorIs(err, ErrTokenExpired)
			suite.Zero(claims)
		})

		suite.Run("невалидный jwt", func() {
			secret := "secret"
			parser := suite.parserWithoutRegistry(secret)

			token := `adsafs.afsfsa.gsdsddsggd`
			claims, err := parser.Parse(token)
			suite.Error(err)
			suite.Zero(claims)
		})

		suite.Run("пустой jwt", func() {
			secret := "secret"
			parser := suite.parserWithoutRegistry(secret)

			claims, err := parser.Parse("")
			suite.Error(err)
			suite.Zero(claims)
		})
	})
	suite.Run("Парсер c Registry", func() {
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
			suite.Parser.Registry = nil
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

			suite.registryMock.EXPECT().IssueTime(sid).Return(time.Time{}, nil)

			claims, err := suite.Parser.Parse(token)
			suite.Require().NoError(err)

			suite.Equal(uid.String(), claims.UserID)
			suite.Equal(sid.String(), claims.SessionID)
		})

		suite.Run("если токен создан после даты анулирования, то токен действителен", func() {
			var (
				uid = uuid.New()
				sid = uuid.New()
			)

			suite.registryMock.EXPECT().IssueTime(sid).Return(time.Now().Add(-time.Hour), nil)

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
			suite.registryMock.EXPECT().IssueTime(sid).Return(time.Now(), nil)

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

			timeNow := time.Now()

			time.Sleep(time.Second)
			token2 := suite.createJWT(secret, map[string]any{
				"UserID":    uid,
				"SessionID": sid,
				"iat":       time.Now().Unix(),
			})

			suite.registryMock.EXPECT().IssueTime(sid).Return(timeNow, nil)
			claims, err := suite.Parser.Parse(token1)
			suite.Require().ErrorIs(err, ErrTokenRevoked)
			suite.Zero(claims)

			claims, err = suite.Parser.Parse(token2)
			suite.Require().NoError(err)
			suite.Equal(sid.String(), claims.SessionID)
			suite.Equal(uid.String(), claims.UserID)
		})
	})
}
