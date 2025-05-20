package service

import (
	"fmt"

	"github.com/markbates/goth/gothic"

	"github.com/saime-0/nice-pea-chat/internal/adapter"
	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type OAuth struct {
	Google adapter.OAuthGoogle
}

type OAuthRegistrationCallbackInput struct {
	Provider string
}

func (o *OAuth) RegistrationCallback(in OAuthRegistrationCallbackInput) (any, error) {

}

type GoogleRegistrationInput struct {
	Code string
}

type GoogleRegistrationOut struct {
	User    domain.User
	Session domain.Session
}

func (o *OAuth) GoogleRegistration(in GoogleRegistrationInput) (GoogleRegistrationOut, error) {
	googleUser, err := o.Google.User(in.Code)
	if err != nil {
		return GoogleRegistrationOut{}, err
	}

	return GoogleRegistrationOut{
		User:    domain.User{},
		Session: domain.Session{},
	}, nil
}

type OAuthRegistrationInitInput struct {
	Provider string
}

type OAuthRegistrationInitOut struct {
	RedirectURL string
}

func (o *OAuth) RegistrationInit(in OAuthRegistrationInitInput) (OAuthRegistrationInitOut, error) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
}
