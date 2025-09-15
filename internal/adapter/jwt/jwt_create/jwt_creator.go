package jwtСreate

import (
	"github.com/cristalhq/jwt/v5"
)

type Issuer struct {
	Secret []byte
}

// создает jwt на основе claims
func (c *Issuer) Issue(claims map[string]any) (string, error) {
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
