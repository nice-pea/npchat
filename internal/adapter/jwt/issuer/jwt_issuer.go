package jwtIssuer

import (
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"

	jwt2 "github.com/nice-pea/npchat/internal/adapter/jwt"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
)

type Issuer struct {
	Config jwt2.Config
}

type customClaims struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	jwt.RegisteredClaims
}

// Issue создает jwt на основе некоторых данных из session
func (c *Issuer) Issue(session sessionn.Session) (string, error) {
	// создаем claims
	claims := customClaims{
		UserID:    session.UserID,
		SessionID: session.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(2 * time.Minute)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
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
