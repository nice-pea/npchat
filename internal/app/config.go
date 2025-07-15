package app

import (
	"github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/controller/http2"
	pgsqlRepository "github.com/nice-pea/npchat/internal/repository/pgsql_repository"
)

type Config struct {
	Pgsql       pgsqlRepository.Config
	Http2       http2.Config
	LogLevel    string
	OAuthGoogle oauth_provider.GoogleConfig
	OAuthGitHub oauth_provider.GitHubConfig
}
