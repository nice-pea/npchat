package service

import (
	"fmt"

	"github.com/markbates/goth/gothic"

	"github.com/saime-0/nice-pea-chat/internal/adapter"
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

func (o *OAuth) GoogleRegistration(in GoogleRegistrationInput) (any, error) {

}

type OAuthRegistrationInitInput struct {
	Provider string
}

type OAuthRegistrationInitOut struct {
	Provider string
}

func (o *OAuth) RegistrationInit(in OAuthRegistrationInitInput) (OAuthRegistrationInitOut, error) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
}
