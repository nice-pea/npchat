package app

import (
	"github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	pgsqlRepository "github.com/nice-pea/npchat/internal/repository/pgsql_repository"
)

type Config struct {
	Pgsql       pgsqlRepository.Config
	HttpAddr    string
	LogLevel    string
	OAuthGoogle oauth_provider.GoogleConfig
	OAuthGitHub oauth_provider.GitHubConfig
}
