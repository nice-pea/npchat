package adapter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type OAuthGoogle interface {
	Exchange(code string) (domain.OAuthToken, error)
	User(token domain.OAuthToken) (domain.OAuthGoogleUser, error)
	AuthCodeURL(state string) string
}

type OAuthGoogleBase struct {
	Config *oauth2.Config
}

func (o *OAuthGoogleBase) Exchange(code string) (domain.OAuthToken, error) {
	token, err := o.Config.Exchange(context.Background(), code)
	if err != nil {
		return domain.OAuthToken{}, err
	}

	return domain.OAuthToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}, nil
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func (o *OAuthGoogleBase) User(token domain.OAuthToken) (domain.OAuthGoogleUser, error) {
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return domain.OAuthGoogleUser{}, err
	}

	defer func() { _ = response.Body.Close() }()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.OAuthGoogleUser{}, err
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
		return domain.OAuthGoogleUser{}, err
	}

	return domain.OAuthGoogleUser(googleUser), nil
}

func (o *OAuthGoogleBase) AuthCodeURL(state string) string {
	return o.Config.AuthCodeURL(state)
}

// MockOAuthGoogle — мок для OAuthGoogle
type MockOAuthGoogle struct {
	ExchangeFunc    func(code string) (domain.OAuthToken, error)
	UserFunc        func(token domain.OAuthToken) (domain.OAuthGoogleUser, error)
	AuthCodeURLFunc func(state string) string
}

func (m *MockOAuthGoogle) Exchange(code string) (domain.OAuthToken, error) {
	if m.ExchangeFunc != nil {
		return m.ExchangeFunc(code)
	}
	panic("Exchange not mocked")
}

func (m *MockOAuthGoogle) User(token domain.OAuthToken) (domain.OAuthGoogleUser, error) {
	if m.UserFunc != nil {
		return m.UserFunc(token)
	}
	panic("User not mocked")
}

func (m *MockOAuthGoogle) AuthCodeURL(state string) string {
	if m.AuthCodeURLFunc != nil {
		return m.AuthCodeURLFunc(state)
	}
	panic("AuthCodeURL not mocked")
}
