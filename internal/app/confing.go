package app

import (
	"github.com/saime-0/nice-pea-chat/internal/adapter/oauth_provider"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
)

type Config struct {
	SQLite      sqlite.Config
	HttpAddr    string
	OAuthGoogle oauth_provider.GoogleConfig
	OAuthGitHub oauth_provider.GitHubConfig
}
