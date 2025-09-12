package jwt

import (
	"time"

	"github.com/cristalhq/jwt/v5"
)

type JWTCreator interface {
	Create(uid string, sid string) string
}

type JWTC struct {
	secret string
	ttl    time.Duration
}

func NewJWTCreator(secret string, ttl time.Duration) JWTC {
	return JWTC{
		secret,
		ttl,
	}
}

// TODO:
// думаю можно сделать чтобы в Create передавали CustomClaims (это упростит тетсирование)
type CustomClaims struct {
	UserID    string
	SessionID string
	jwt.RegisteredClaims
}

func (c *JWTC) Create(uid string, sid string) (string, error) {

	// create a Signer (HMAC in this example)
	key := []byte(c.secret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return "", err
	}
	// create claims (you can create your own, see: ExampleBuilder_withUserClaims)
	claims := CustomClaims{
		UserID:    uid,
		SessionID: sid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "npchat",
			Subject:   "authentication",
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(c.ttl)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
		},
	}

	// create a Builder
	builder := jwt.NewBuilder(signer)

	// and build a Token
	token, err := builder.Build(claims)
	if err != nil {
		return "", err
	}

	return token.String(), nil
}
