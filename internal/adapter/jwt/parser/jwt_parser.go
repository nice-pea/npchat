package jwtParser

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"

	jwt2 "github.com/nice-pea/npchat/internal/adapter/jwt"

	redisCache "github.com/nice-pea/npchat/internal/adapter/jwt/repository/redis"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
)

type OutJWT struct {
	UserID    string
	SessionID string
}

type Parser struct {
	Config jwt2.Config
	cache  redisCache.JWTIssuanceRegistry
}

var (
	ErrTimeOut      = errors.New("время жизни токена истекло")
	ErrTokenRevoked = errors.New("токен аннулирован")
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

func (p *Parser) ValidateJWTWithInvalidation(token string) (middleware.OutJwt, error) {

	// create a Verifier (HMAC in this example)
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(p.Config.SecretKey))

	if err != nil {
		return middleware.OutJwt{}, err
	}

	// parse and verify a token
	tokenBytes := []byte(token)
	newToken, err := jwt.Parse(tokenBytes, verifier)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	err = verifier.Verify(newToken)
	if err != nil {
		return middleware.OutJwt{}, err
	}
	// get Registered claims
	var newClaims CustomClaims
	errClaims := json.Unmarshal(newToken.Claims(), &newClaims)
	if errClaims != nil {
		return middleware.OutJwt{}, err
	}
	// verify claims as you wish
	if !newClaims.IsValidAt(time.Now()) {
		return middleware.OutJwt{}, ErrTimeOut
	}

	sessionId, err := uuid.Parse(newClaims.SessionID)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	issuedAt := newClaims.IssuedAt

	timefromCache, err := p.cache.GetIssueTime(sessionId)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	if timefromCache.After(issuedAt.Time) {
		return middleware.OutJwt{}, ErrTokenRevoked
	}

	return customClaimsToOutJWT(newClaims), nil

}
