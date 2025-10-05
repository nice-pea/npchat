package jwtParser

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"

	jwt2 "github.com/nice-pea/npchat/internal/adapter/jwt"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
)

// Register интерфейс для репозитория хранения и полчения дат анулирования
type Registry interface {
	RegisterIssueTime(sessionID uuid.UUID, issueTime time.Time) error
	IssueTime(sessionID uuid.UUID) (time.Time, error)
}

// OutJWT структура данных JWT
type OutJWT struct {
	UserID    string
	SessionID string
}

// Parser парсер JWT
type Parser struct {
	Config   jwt2.Config
	Registry Registry
}

// Ошибки модуля
var (
	ErrIatEmpty     = errors.New("время создания токена отсутствует")
	ErrTokenExpired = errors.New("время жизни токена истекло")
	ErrTokenRevoked = errors.New("токен аннулирован")
)

// CustomClaims структура для хранения данных из claims
type CustomClaims struct {
	UserID    string
	SessionID string
	jwt.RegisteredClaims
}

// customClaimsToOutJWT преобразовывает CustomClaims в middleware.OutJwt
func customClaimsToOutJWT(cc CustomClaims) middleware.OutJwt {
	return middleware.OutJwt{
		UserID:    cc.UserID,
		SessionID: cc.SessionID,
	}
}

// getClaims парсит токен в CustomClaims
func (p *Parser) getClaims(token string) (CustomClaims, error) {
	// Создание верификатора с секретным ключом
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(p.Config.SecretKey))

	if err != nil {
		return CustomClaims{}, err
	}

	// Парсинг токена
	newToken, err := jwt.Parse([]byte(token), verifier)
	if err != nil {
		return CustomClaims{}, err
	}
	// Подпись токена
	if err = verifier.Verify(newToken); err != nil {
		return CustomClaims{}, err
	}

	// Распаковка claims
	var newClaims CustomClaims
	errClaims := json.Unmarshal(newToken.Claims(), &newClaims)
	if errClaims != nil {
		return CustomClaims{}, errClaims
	}

	// Проверка временных ограничений
	if !newClaims.IsValidAt(time.Now()) {
		return CustomClaims{}, ErrTokenExpired
	}
	return newClaims, nil
}

// parse парсит токен без проверки анулирования
func (p *Parser) parse(token string) (middleware.OutJwt, error) {
	claims, err := p.getClaims(token)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	return customClaimsToOutJWT(claims), nil
}

// parseAndValidateJWTWithInvalidation парсит токен и проверяет его на аннулирование
func (p *Parser) parseAndValidateJWTWithInvalidation(token string) (middleware.OutJwt, error) {
	claims, err := p.getClaims(token)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	sessionId, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	// Получение времени последнего выпуска токена из Redis
	timefromCache, err := p.Registry.IssueTime(sessionId)
	if err != nil {
		return middleware.OutJwt{}, err
	}

	issuedAt := claims.IssuedAt
	if issuedAt == nil {
		return middleware.OutJwt{}, ErrIatEmpty
	}

	// Если время в кэше больше, чем в токене - токен считается аннулированным
	if timefromCache.After(issuedAt.Time) {
		return middleware.OutJwt{}, ErrTokenRevoked
	}

	return customClaimsToOutJWT(claims), nil
}

// Parse парсит токен
// Выбирает режим проверки в зависимости от конфигурации
func (p *Parser) Parse(token string) (middleware.OutJwt, error) {
	if p.Config.VerifyTokenWithInvalidation && p.Registry != nil {
		return p.parseAndValidateJWTWithInvalidation(token)
	}

	return p.parse(token)
}
