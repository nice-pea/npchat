package app

import (
	"github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/repository/sqlite"
)

type Config struct {
	SQLite      sqlite.Config
	HttpAddr    string
	OAuthGoogle oauth_provider.GoogleConfig
	OAuthGitHub oauth_provider.GitHubConfig
}
