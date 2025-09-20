package app

import (
	jwt2 "github.com/nice-pea/npchat/internal/adapter/jwt"
	oauthProvider "github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/controller/http2"
	pgsqlRepository "github.com/nice-pea/npchat/internal/repository/pgsql_repository"
)

type Config struct {
	Pgsql       pgsqlRepository.Config
	Http2       http2.Config
	LogLevel    string
	OauthGoogle oauthProvider.GoogleConfig
	OauthGithub oauthProvider.GithubConfig
	Jwt         jwt2.Config
}
