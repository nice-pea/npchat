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

<<<<<<< HEAD:internal/adapter/jwt/parser/jwt_parser.go
// Parse разбирает токен и возвращает данные из него
func (p *Parser) Parse(token string) (middleware.OutJwt, error) {
	// Создать валидатор
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(p.Config.SecretKey))
=======
func (p *JWTParser) getClaims(token string) (CustomClaims, error) {
	// create a Verifier (HMAC in this example)

	verifier, err := jwt.NewVerifierHS(jwt.HS256, p.Secret)

>>>>>>> c6ff310 (вынести код в отдельную функцию):internal/adapter/jwt/jwt_parse/jwt_parser.go
	if err != nil {
		return CustomClaims{}, err
	}

	// Разобрать токен и проверить его
	newToken, err := jwt.Parse([]byte(token), verifier)
	if err != nil {
		return CustomClaims{}, err
	}
	if err = verifier.Verify(newToken); err != nil {
		return middleware.OutJwt{}, err
	}

<<<<<<< HEAD:internal/adapter/jwt/parser/jwt_parser.go
	// Получить данные из токена
=======
	err = verifier.Verify(newToken)
	if err != nil {
		return CustomClaims{}, err
	}
	// get Registered claims
>>>>>>> c6ff310 (вынести код в отдельную функцию):internal/adapter/jwt/jwt_parse/jwt_parser.go
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

<<<<<<< HEAD:internal/adapter/jwt/parser/jwt_parser.go
func (p *Parser) ValidateJWTWithInvalidation(token string) (middleware.OutJwt, error) {

	// create a Verifier (HMAC in this example)
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(p.Config.SecretKey))

=======
func (p *JWTParser) Parse(token string) (middleware.OutJwt, error) {
	claims, err := p.getClaims(token)
>>>>>>> c6ff310 (вынести код в отдельную функцию):internal/adapter/jwt/jwt_parse/jwt_parser.go
	if err != nil {
		return middleware.OutJwt{}, err
	}

	return customClaimsToOutJWT(claims), nil
}

func (p *JWTParser) ParseAndValidateJWTWithInvalidation(token string) (middleware.OutJwt, error) {
	claims, err := p.getClaims(token)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	sessionId, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	issuedAt := claims.IssuedAt

	timefromCache, err := p.cache.GetIssueTime(sessionId)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	if timefromCache.After(issuedAt.Time) {
		return middleware.OutJwt{}, ErrTokenRevoked
	}

	return customClaimsToOutJWT(claims), nil

}
