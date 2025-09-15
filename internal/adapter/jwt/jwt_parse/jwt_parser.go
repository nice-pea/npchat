package jwtParser

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/cristalhq/jwt/v5"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
)

type OutJWT struct {
	UserID    string
	SessionID string
}

type JWTParser struct {
	Secret []byte
}

var (
	ErrTimeOut = errors.New("время жизни токена истекло")
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

func (p *JWTParser) Parse(token string) (middleware.OutJwt, error) {
	// create a Verifier (HMAC in this example)

	verifier, err := jwt.NewVerifierHS(jwt.HS256, p.Secret)

	if err != nil {
		return middleware.OutJwt{}, err
	}

	// parse and verify a token
	tokenBytes := []byte(token)
	newToken, err := jwt.Parse(tokenBytes, verifier)
	if err != nil {
		return middleware.OutJwt{}, err
	}
	// or just verify it's signature
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
	return customClaimsToOutJWT(newClaims), nil
}
