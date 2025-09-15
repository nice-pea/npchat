package app

import (
	"log/slog"

	eventsBus "github.com/nice-pea/npchat/internal/adapter/events_bus"
	jwtIssuer "github.com/nice-pea/npchat/internal/adapter/jwt/jwt_create"
	jwtParser "github.com/nice-pea/npchat/internal/adapter/jwt/jwt_parse"
	oauthProvider "github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
)

type adapters struct {
	oauthProviders oauth.Providers
	eventBus       *eventsBus.EventsBus
	jwtUtils       JwtUtils
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
		jwtUtils:       initJwtUtils(cfg.JwtSecret),
	}
}

type JwtUtils struct {
	*jwtParser.JWTParser
	*jwtIssuer.Issuer
}

func initJwtUtils(secret string) JwtUtils {

	return JwtUtils{
		JWTParser: &jwtParser.JWTParser{Secret: []byte(secret)},
		Issuer:    &jwtIssuer.Issuer{Secret: []byte(secret)},
	}
}
