package app

import (
	"log/slog"

	eventsBus "github.com/nice-pea/npchat/internal/adapter/events_bus"
	oauthProvider "github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
)

type adapters struct {
	oauthProviders oauth.OAuthProviders
	eventBus       *eventsBus.EventsBus
}

func (a *adapters) OauthProviders() oauth.OauthProviders {
	return a.oauthProviders
}

func initAdapters(cfg Config) *adapters {
	oauthProviders := oauth.OauthProviders{}
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
