package adapter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type OAuthGoogle interface {
	Exchange(code string) (GoogleToken, error)
	User(code string) (GoogleUser, error)
	AuthCodeURL(state string) string
}

type GoogleUser struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
	Locale        string
}

type OAuthGoogleBase struct {
	Config *oauth2.Config
}

type GoogleToken struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string `json:"access_token"`

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string `json:"token_type,omitempty"`

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string `json:"refresh_token,omitempty"`

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time `json:"expiry,omitempty"`
}

func (o *OAuthGoogleBase) Exchange(code string) (GoogleToken, error) {
	token, err := o.Config.Exchange(context.Background(), code)
	if err != nil {
		return GoogleToken{}, err
	}

	return GoogleToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}, nil
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func (o *OAuthGoogleBase) User(code string) (GoogleUser, error) {
	token, err := o.Exchange(code)
	if err != nil {
		return GoogleUser{}, err
	}

	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return GoogleUser{}, err
	}

	defer func() { _ = response.Body.Close() }()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return GoogleUser{}, err
	}

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}
	if err = json.Unmarshal(data, &googleUser); err != nil {
		return GoogleUser{}, err
	}

	return GoogleUser(googleUser), nil
}

func (o *OAuthGoogleBase) AuthCodeURL(state string) string {
	return o.Config.AuthCodeURL(state)
}
