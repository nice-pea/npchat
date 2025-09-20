package jwtParser

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/cristalhq/jwt/v5"

	jwt2 "github.com/nice-pea/npchat/internal/adapter/jwt"
	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
)

type OutJWT struct {
	UserID    string
	SessionID string
}

type Parser struct {
	Config jwt2.Config
}

var (
	ErrTimeOut = errors.New("время жизни токена истекло")
)

type CustomClaims struct {
	UserID    string
	SessionID string
	jwt.RegisteredClaims
}

func customClaimsToOutJWT(cc CustomClaims) middleware.OutJwt {
	return middleware.OutJwt{
		UserID:    cc.UserID,
		SessionID: cc.SessionID,
	}
}

// Parse разбирает токен и возвращает данные из него
func (p *Parser) Parse(token string) (middleware.OutJwt, error) {
	// Создать валидатор
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(p.Config.SecretKey))
	if err != nil {
		return middleware.OutJwt{}, err
	}

	// Разобрать токен и проверить его
	newToken, err := jwt.Parse([]byte(token), verifier)
	if err != nil {
		return middleware.OutJwt{}, err
	}
	if err = verifier.Verify(newToken); err != nil {
		return middleware.OutJwt{}, err
	}

	// Получить данные из токена
	var newClaims CustomClaims
	errClaims := json.Unmarshal(newToken.Claims(), &newClaims)
	if errClaims != nil {
		return middleware.OutJwt{}, err
	}

	// Валидация времени жизни токена
	if !newClaims.IsValidAt(time.Now()) {
		return middleware.OutJwt{}, ErrTimeOut
	}
	return customClaimsToOutJWT(newClaims), nil
}
