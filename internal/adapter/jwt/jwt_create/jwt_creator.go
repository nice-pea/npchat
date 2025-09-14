package jwtCreator

import (
	"github.com/cristalhq/jwt/v5"
)

type JwtIssuer interface {
	Issue(claims map[string]any) (string, error)
}

type JWTC struct {
	secret string
}

func NewJWTCreator(secret string) JWTC {
	return JWTC{
		secret,
	}
}

// создает jwt на основе claims
func (c *JWTC) Issue(claims map[string]any) (string, error) {
	// создаем Signer
	key := []byte(c.secret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
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
