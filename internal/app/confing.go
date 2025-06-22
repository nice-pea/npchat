package app

import (
	"log/slog"

	"github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/repository/sqlite"
)

type Config struct {
	SQLite      sqlite.Config
	HttpAddr    string
	SlogLevel   slog.Level
	OAuthGoogle oauth_provider.GoogleConfig
	OAuthGitHub oauth_provider.GitHubConfig
}
