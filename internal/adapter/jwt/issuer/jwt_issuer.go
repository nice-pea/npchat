package jwtIssuer

import (
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"

	jwt2 "github.com/nice-pea/npchat/internal/adapter/jwt"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
)

// Issuer генератор JWT-токенов
type Issuer struct {
	Config jwt2.Config
}

// customClaims расширенные claims для JWT
type customClaims struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	jwt.RegisteredClaims
}

// Issue создает JWT-токен на основе данных сессии
func (c *Issuer) Issue(session sessionn.Session) (string, error) {
	// Получение текущего времени
	nowTime := time.Now()
	
	// Формирование claims токена
	claims := customClaims{
		UserID:    session.UserID,
		SessionID: session.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: nowTime.Add(2 * time.Minute)},
			IssuedAt:  &jwt.NumericDate{Time: nowTime},
		},
	}

	// Создать подпись
	signer, err := jwt.NewSignerHS(jwt.HS256, []byte(c.Config.SecretKey))
	if err != nil {
		return "", err
	}

	// Создать токен
	token, err := jwt.NewBuilder(signer).Build(claims)
	if err != nil {
		return "", err
	}

	return token.String(), nil
}
