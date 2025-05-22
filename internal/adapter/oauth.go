package adapter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"
	googleOAuth "golang.org/x/oauth2/google"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

const (
	OAuthProviderGoogle = "google"
	OAuthProviderGithub = "github"
)

const (
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	oauthGitHubUrlAPI = "https://api.github.com/user"
)

type OAuthProviders = map[string]OAuthProvider

type OAuthProvider interface {
	Exchange(code string) (domain.OAuthToken, error)
	User(token domain.OAuthToken) (domain.OAuthUser, error)
	AuthCodeURL(state string) string
}

type OAuthGoogle struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (o *OAuthGoogle) config() *oauth2.Config {
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

func (o *OAuthGoogle) Exchange(code string) (domain.OAuthToken, error) {
	token, err := o.config().Exchange(context.Background(), code)
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

func (o *OAuthGoogle) User(token domain.OAuthToken) (domain.OAuthUser, error) {
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
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
		Provider: OAuthProviderGoogle,
	}, nil
}

func (o *OAuthGoogle) AuthCodeURL(state string) string {
	return o.config().AuthCodeURL(state)
}

type OAuthGitHub struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (o *OAuthGitHub) config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		Endpoint:     githubOAuth.Endpoint,
		RedirectURL:  o.RedirectURL,
		Scopes:       []string{"user:email"},
	}
}

func (o *OAuthGitHub) Exchange(code string) (domain.OAuthToken, error) {
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
		Provider:     OAuthProviderGithub,
	}, nil
}

func (o *OAuthGitHub) User(token domain.OAuthToken) (domain.OAuthUser, error) {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.AccessToken},
	))

	response, err := client.Get(oauthGitHubUrlAPI)
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
		Provider: OAuthProviderGithub,
	}, nil
}

func (o *OAuthGitHub) AuthCodeURL(state string) string {
	return o.config().AuthCodeURL(state)
}

//type OAuthGithub struct {
//	Config *oauth2.Config
//}
//
//func (o *OAuthGithub) Exchange(code string) (domain.OAuthToken, error) {
//	token, err := o.Config.Exchange(context.Background(), code)
//	if err != nil {
//		return domain.OAuthToken{}, err
//	}
//
//	return domain.OAuthToken{
//		AccessToken:  token.AccessToken,
//		TokenType:    token.TokenType,
//		RefreshToken: token.RefreshToken,
//		Expiry:       token.Expiry,
//	}, nil
//}
//
//func (o *OAuthGithub) User(token domain.OAuthToken) (domain.OAuthUser, error) {
//	client := &http.Client{}
//	req, err := http.NewRequest("GET", oauthGithubUrlAPI, nil)
//	if err != nil {
//		return domain.OAuthUser{}, err
//	}
//
//	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
//	response, err := client.Do(req)
//	if err != nil {
//		return domain.OAuthUser{}, err
//	}
//
//	defer func() { _ = response.Body.Close() }()
//	data, err := io.ReadAll(response.Body)
//	if err != nil {
//		return domain.OAuthUser{}, err
//	}
//
//	var githubUser struct {
//		ID    string `json:"id"`
//		Email string `json:"email"`
//		Name  string `json:"name"`
//	}
//	if err = json.Unmarshal(data, &githubUser); err != nil {
//		return domain.OAuthUser{}, err
//	}
//
//	return domain.OAuthUser{
//		ID:       githubUser.ID,
//		Email:    githubUser.Email,
//		Name:     githubUser.Name,
//		Picture:  "",
//		Provider: "",
//	}, nil
//}
//
//func (o *OAuthGithub) AuthCodeURL(state string) string {
//	return o.Config.AuthCodeURL(state)
//}

// OAuthMock мок для OAuthProvider
type OAuthMock struct {
	ExchangeFunc    func(code string) (domain.OAuthToken, error)
	UserFunc        func(token domain.OAuthToken) (domain.OAuthUser, error)
	AuthCodeURLFunc func(state string) string
}

func (m *OAuthMock) Exchange(code string) (domain.OAuthToken, error) {
	if m.ExchangeFunc != nil {
		return m.ExchangeFunc(code)
	}
	panic("Exchange not mocked")
}

func (m *OAuthMock) User(token domain.OAuthToken) (domain.OAuthUser, error) {
	if m.UserFunc != nil {
		return m.UserFunc(token)
	}
	panic("User not mocked")
}

func (m *OAuthMock) AuthCodeURL(state string) string {
	if m.AuthCodeURLFunc != nil {
		return m.AuthCodeURLFunc(state)
	}
	panic("AuthCodeURL not mocked")
}
