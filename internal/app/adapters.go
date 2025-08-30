package app

import (
	"log/slog"

	"github.com/nice-pea/npchat/internal/adapter"
	oauthProvider "github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/service/users/oauth"
)

type adapters struct {
	oauthProviders oauth.OAuthProviders
	discovery      adapter.ServiceDiscovery
}

func (a *adapters) Discovery() adapter.ServiceDiscovery {
	return a.discovery
}

func (a *adapters) OAuthProviders() oauth.OAuthProviders {
	return a.oauthProviders
}

func initAdapters(cfg Config) *adapters {
	discovery := &adapter.ServiceDiscoveryBase{
		Debug: true,
	}
	oauthProviders := oauth.OAuthProviders{}
	if cfg.OAuthGoogle != (oauthProvider.GoogleConfig{}) {
		oauthProviders.Add(oauthProvider.NewGoogle(cfg.OAuthGoogle))
		slog.Info("Подключен OAuth провайдер Google")
	}
	if cfg.OAuthGitHub != (oauthProvider.GitHubConfig{}) {
		oauthProviders.Add(oauthProvider.NewGitHub(cfg.OAuthGitHub))
		slog.Info("Подключен OAuth провайдер GitHub")
	}

	return &adapters{
		oauthProviders: oauthProviders,
		discovery:      discovery,
	}
}
