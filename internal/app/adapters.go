package app

import (
	"log/slog"

	eventsBus "github.com/nice-pea/npchat/internal/adapter/events_bus"
	jwt2 "github.com/nice-pea/npchat/internal/adapter/jwt"
	jwtIssuer "github.com/nice-pea/npchat/internal/adapter/jwt/issuer"
	jwtParser "github.com/nice-pea/npchat/internal/adapter/jwt/parser"
	oauthProvider "github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	registerHandler "github.com/nice-pea/npchat/internal/controller/http2/register_handler"
	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
)

type adapters struct {
	oauthProviders oauth.Providers
	eventBus       *eventsBus.EventsBus
	jwtParser      middleware.JwtParser
	jwtIssuer      registerHandler.JwtIssuer
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

	// Включить jwt аутентификацию, если конфиг jwt задан
	var jwtParser2 middleware.JwtParser
	var jwtIssuer2 registerHandler.JwtIssuer
	if cfg.Jwt != (jwt2.Config{}) {
		jwtParser2 = &jwtParser.Parser{Config: cfg.Jwt}
		jwtIssuer2 = &jwtIssuer.Issuer{Config: cfg.Jwt}
		slog.Info("Подключена jwt аутентификация")
	}

	return &adapters{
		oauthProviders: oauthProviders,
		eventBus:       new(eventsBus.EventsBus),
		jwtParser:      jwtParser2,
		jwtIssuer:      jwtIssuer2,
	}
}
