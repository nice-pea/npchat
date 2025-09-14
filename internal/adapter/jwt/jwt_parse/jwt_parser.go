package jwt_parse

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
	secret string
}

func NewJWTParser(secret string) JWTParser {
	return JWTParser{secret}
}

var (
	ErrTimeOut = errors.New("время жизни токена истекло")
)

type CustomClaims struct {
	UserID    string
	SessionID string
	jwt.RegisteredClaims
}

func customClaimsToOutJWT(cc CustomClaims) middleware.OutJWT {
	return middleware.OutJWT{
		UserID:    cc.UserID,
		SessionID: cc.SessionID,
	}
}

func (p *JWTParser) Parse(token string) (middleware.OutJWT, error) {
	// create a Verifier (HMAC in this example)
	key := []byte(p.secret)
	verifier, err := jwt.NewVerifierHS(jwt.HS256, key)

	if err != nil {
		return middleware.OutJWT{}, err
	}

	// parse and verify a token
	tokenBytes := []byte(token)
	newToken, err := jwt.Parse(tokenBytes, verifier)
	if err != nil {
		return middleware.OutJWT{}, err
	}
	// or just verify it's signature
	err = verifier.Verify(newToken)
	if err != nil {
		return middleware.OutJWT{}, err
	}
	// get Registered claims
	var newClaims CustomClaims
	errClaims := json.Unmarshal(newToken.Claims(), &newClaims)
	if errClaims != nil {
		return middleware.OutJWT{}, err
	}
	// verify claims as you wish
	var isValidTime bool = newClaims.IsValidAt(time.Now())
	if !isValidTime {
		return middleware.OutJWT{}, ErrTimeOut
	}
	return customClaimsToOutJWT(newClaims), nil
}
