package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/saime-0/nice-pea-chat/internal/adapter"
	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type OAuth struct {
	Google adapter.OAuthGoogle
}

type OAuthRegistrationCallbackInput struct {
	Provider string
}

type GoogleRegistrationInput struct {
	Code  string
	State string
}

type GoogleRegistrationOut struct {
	User    domain.User
	Session domain.Session
}

func (o *OAuth) GoogleRegistration(in GoogleRegistrationInput) (GoogleRegistrationOut, error) {
	if err := validateStateOauthJWT(in.State, JWTSecret); err != nil {
		return GoogleRegistrationOut{}, err
	}

	googleUser, err := o.Google.User(in.Code)
	if err != nil {
		return GoogleRegistrationOut{}, err
	}

	fmt.Println(googleUser)

	return GoogleRegistrationOut{
		User:    domain.User{},
		Session: domain.Session{},
	}, nil
}

//type GoogleRegistrationInitInput struct {
//	RedirectURL string
//}

type GoogleRegistrationInitOut struct {
	RedirectURL string
}

var JWTSecret = os.Getenv("OAUTH_JWT_SECRET")

func (o *OAuth) GoogleRegistrationInit() (GoogleRegistrationInitOut, error) {
	oauthState, err := generateStateOauthJWT(JWTSecret)
	if err != nil {
		return GoogleRegistrationInitOut{}, err
	}

	return GoogleRegistrationInitOut{
		RedirectURL: o.Google.AuthCodeURL(oauthState),
	}, nil
}

func generateStateOauthJWT(secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(10 * time.Minute).Unix(), // State живёт 10 минут
	})
	return token.SignedString([]byte(secret))
}

func validateStateOauthJWT(stateToken, secret string) error {
	token, err := jwt.Parse(stateToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token is not valid")
	}

	return nil
}
