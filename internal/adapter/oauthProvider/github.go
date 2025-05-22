package oauthProvider

import (
	"context"
	"encoding/json"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type GitHub struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (o *GitHub) Name() string {
	return "github"
}

func (o *GitHub) config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  o.RedirectURL,
		Scopes:       []string{"user:email"},
	}
}

func (o *GitHub) Exchange(code string) (domain.OAuthToken, error) {
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

func (o *GitHub) User(token domain.OAuthToken) (domain.OAuthUser, error) {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.AccessToken},
	))

	const getUser = "https://api.github.com/user"
	response, err := client.Get(getUser)
	if err != nil {
		return domain.OAuthUser{}, err
	}
	defer func() { _ = response.Body.Close() }()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.OAuthUser{}, err
	}

	var githubUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.Unmarshal(data, &githubUser); err != nil {
		return domain.OAuthUser{}, err
	}

	return domain.OAuthUser{
		ID:       string(rune(githubUser.ID)),
		Email:    githubUser.Email,
		Name:     githubUser.Name,
		Picture:  githubUser.AvatarURL,
		Provider: o.Name(),
	}, nil
}

func (o *GitHub) AuthCodeURL(state string) string {
	return o.config().AuthCodeURL(state)
}
