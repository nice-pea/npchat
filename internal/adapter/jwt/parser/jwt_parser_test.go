package jwtParser

import (
	"testing"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jwt2 "github.com/nice-pea/npchat/internal/adapter/jwt"
)

func Test_JWTParser_Parse(t *testing.T) {
	t.Run("валидный jwt можно разобрать и получить данные", func(t *testing.T) {
		secret := "secret"
		parser := Parser{jwt2.Config{SecretKey: secret}}

		var (
			uid = "123"
			sid = "456"
		)

		token, err := createJWT(secret, map[string]any{
			"UserID":    uid,
			"SessionID": sid,
		})
		require.NoError(t, err)

		claims, err := parser.Parse(token)
		require.NoError(t, err)

		assert.Equal(t, uid, claims.UserID)
		assert.Equal(t, sid, claims.SessionID)
	})
	t.Run("jwt существующий больше exp - невалиден", func(t *testing.T) {
		secret := "secret"
		parser := Parser{jwt2.Config{SecretKey: secret}}

		token, err := createJWT(secret, map[string]any{
			"UserID":    "123",
			"SessionID": "456",
			"exp":       time.Now().Add(time.Millisecond).Unix(),
		})
		require.NoError(t, err)
		time.Sleep(2 * time.Millisecond)

		claims, err := parser.Parse(token)
		assert.ErrorIs(t, err, ErrTimeOut)
		assert.Zero(t, claims)
	})
	t.Run("jwt существующий меньше exp - валиден", func(t *testing.T) {
		secret := "secret"
		parser := Parser{jwt2.Config{SecretKey: secret}}

		token, err := createJWT(secret, map[string]any{
			"UserID":    "123",
			"SessionID": "456",
			"exp":       time.Now().Add(1000 * time.Millisecond).Unix(),
		})
		require.NoError(t, err)
		time.Sleep(100 * time.Millisecond)

		claims, err := parser.Parse(token)
		require.NoError(t, err)
		assert.NotZero(t, claims)
	})
	t.Run("jwt существующий больше 2 минут - невалиден", func(t *testing.T) {
		secret := "secret"
		parser := Parser{jwt2.Config{SecretKey: secret}}

		//token - истекший jwt токен
		token := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIxMjMiLCJTZXNzaW9uSUQiOiI0NTYiLCJpc3MiOiJucGNoYXQiLCJzdWIiOiJhdXRoZW50aWNhdGlvbiIsImV4cCI6MTc1NzcwMDM2NCwiaWF0IjoxNzU3NzAwMjQ0LCJuYmYiOjE3NTc3MDAyNDR9.kpKiS63GV1XQYapTC9jxAlACoKToOIWISgzvJIVeZ2I`

		claims, err := parser.Parse(token)
		assert.ErrorIs(t, err, ErrTimeOut)
		assert.Zero(t, claims)
	})

	t.Run("невалидный jwt", func(t *testing.T) {
		secret := "secret"
		parser := Parser{jwt2.Config{SecretKey: secret}}

		token := `adsafs.afsfsa.gsdsddsggd`

		claims, err := parser.Parse(token)
		assert.Error(t, err)
		assert.Zero(t, claims)
	})
	t.Run("пустой jwt", func(t *testing.T) {
		secret := "secret"
		parser := Parser{jwt2.Config{SecretKey: secret}}

		claims, err := parser.Parse("")
		assert.Error(t, err)
		assert.Zero(t, claims)
	})
}

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
