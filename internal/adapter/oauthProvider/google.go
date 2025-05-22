package oauthProvider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	googleOAuth "golang.org/x/oauth2/google"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type Google struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (o *Google) Name() string {
	return "google"
}

func (o *Google) config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		Endpoint:     googleOAuth.Endpoint,
		RedirectURL:  o.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
}

func (o *Google) Exchange(code string) (domain.OAuthToken, error) {
	token, err := o.config().Exchange(context.Background(), code)
	if err != nil {
		return domain.OAuthToken{}, err
	}

	return domain.OAuthToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		LinkID:       "",
		Provider:     o.Name(),
	}, nil
}

func (o *Google) User(token domain.OAuthToken) (domain.OAuthUser, error) {
	const getUser = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	response, err := http.Get(getUser + token.AccessToken)
	if err != nil {
		return domain.OAuthUser{}, err
	}

	defer func() { _ = response.Body.Close() }()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.OAuthUser{}, err
	}

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err = json.Unmarshal(data, &googleUser); err != nil {
		return domain.OAuthUser{}, err
	}

	return domain.OAuthUser{
		ID:       googleUser.ID,
		Email:    googleUser.Email,
		Name:     googleUser.Name,
		Picture:  googleUser.Picture,
		Provider: o.Name(),
	}, nil
}

func (o *Google) AuthCodeURL(state string) string {
	return o.config().AuthCodeURL(state)
}
