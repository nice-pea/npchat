package jwtIssuer

import (
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
)

type Issuer struct {
	Secret []byte
}

type customClaims struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	jwt.RegisteredClaims
}

// создает jwt на основе некоторых данных из session
func (c *Issuer) Issue(session sessionn.Session) (string, error) {
	// создаем claims
	claims := customClaims{
		UserID:    session.UserID,
		SessionID: session.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(2 * time.Minute)},
		},
	}

	// создаем Signer
	signer, err := jwt.NewSignerHS(jwt.HS256, c.Secret)
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
