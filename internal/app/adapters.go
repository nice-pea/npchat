package app

import (
	"log/slog"

	eventsBus "github.com/nice-pea/npchat/internal/adapter/events_bus"
	"github.com/nice-pea/npchat/internal/adapter"
	"github.com/nice-pea/npchat/internal/adapter/jwt/jwt_create"
	"github.com/nice-pea/npchat/internal/adapter/jwt/jwt_parse"

	oauthProvider "github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
)

type adapters struct {
	oauthProviders oauth.Providers
	eventBus       *eventsBus.EventsBus
}

func (a *adapters) OauthProviders() oauth.Providers {
	return a.oauthProviders
}

func initAdapters(cfg Config) *adapters {
	oauthProviders := oauth.Providers{}
	if cfg.OauthGoogle != (oauthProvider.GoogleConfig{}) {
		oauthProviders.Add(oauthProvider.NewGoogle(cfg.OauthGoogle))
		slog.Info("Подключен Oauth провайдер Google")
	}
	if cfg.OauthGithub != (oauthProvider.GithubConfig{}) {
		oauthProviders.Add(oauthProvider.NewGithub(cfg.OauthGithub))
		slog.Info("Подключен Oauth провайдер Github")
	}

	return &adapters{
		oauthProviders: oauthProviders,
		eventBus:       new(eventsBus.EventsBus),
	}
}

type JWTUtils struct {
	*jwt_parse.JWTParser
	*jwt_create.JWTC
}

func initJwtUtils(secret string) JWTUtils {
	parser := jwt_parse.NewJWTParser(secret)
	issuer := jwt_create.NewJWTCreator(secret)
	return JWTUtils{
		JWTParser: &parser,
		JWTC:      &issuer,
	}
}
