package adapter

type OAuthGoogle interface {
	Exchange(code string)
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type UserInfo struct {
	Email string
	Name  string
}
