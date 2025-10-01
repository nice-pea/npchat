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

	VerifyTokenWithAdvancedChecks bool
	Registry                      redisCache.Registry
}

var (
	ErrIatEmpty     = errors.New("время создания токена отстутсвует")
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
func (p *Parser) getClaims(token string) (CustomClaims, error) {
	// Создать валидатор
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(p.Config.SecretKey))

	if err != nil {
		return CustomClaims{}, err
	}

	// Разобрать токен и проверить его
	newToken, err := jwt.Parse([]byte(token), verifier)
	if err != nil {
		return CustomClaims{}, err
	}
	if err = verifier.Verify(newToken); err != nil {
		return CustomClaims{}, err
	}

	// Получить данные из токена
	var newClaims CustomClaims
	errClaims := json.Unmarshal(newToken.Claims(), &newClaims)
	if errClaims != nil {
		return CustomClaims{}, err
	}

	// Валидация времени жизни токена
	if !newClaims.IsValidAt(time.Now()) {
		return CustomClaims{}, ErrTimeOut
	}
	return newClaims, nil
}

func (p *Parser) parse(token string) (middleware.OutJwt, error) {
	claims, err := p.getClaims(token)

	if err != nil {
		return middleware.OutJwt{}, err
	}

	return customClaimsToOutJWT(claims), nil
}

func (p *Parser) parseAndValidateJWTWithInvalidation(token string) (middleware.OutJwt, error) {
	claims, err := p.getClaims(token)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	sessionId, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	timefromCache, err := p.Registry.IssueTime(sessionId)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	issuedAt := claims.IssuedAt
	if issuedAt == nil {
		return middleware.OutJwt{}, ErrIatEmpty
	}

	if timefromCache.After(issuedAt.Time) {
		return middleware.OutJwt{}, ErrTokenRevoked
	}

	return customClaimsToOutJWT(claims), nil

}

func (p *Parser) Parse(token string) (middleware.OutJwt, error) {
	if p.VerifyTokenWithAdvancedChecks && p.Registry.Cli != nil {
		return p.parseAndValidateJWTWithInvalidation(token)
	}

	return p.parse(token)
}
