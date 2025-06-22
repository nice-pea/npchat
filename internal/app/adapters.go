package app

import (
	"github.com/nice-pea/npchat/internal/adapter"
	oauthProvider "github.com/nice-pea/npchat/internal/adapter/oauth_provider"
	"github.com/nice-pea/npchat/internal/service"
)

type adapters struct {
	oauthProviders service.OAuthProviders
	discovery      adapter.ServiceDiscovery
}

func (a *adapters) Discovery() adapter.ServiceDiscovery {
	return a.discovery
}

func (a *adapters) OAuthProviders() service.OAuthProviders {
	return a.oauthProviders
}

func initAdapters(cfg Config) *adapters {
	discovery := &adapter.ServiceDiscoveryBase{
		Debug: true,
	}
	oauthProviders := service.OAuthProviders{}
	if cfg.OAuthGoogle != (oauthProvider.GoogleConfig{}) {
		oauthProviders.Add(oauthProvider.NewGoogle(cfg.OAuthGoogle))
	}
	if cfg.OAuthGitHub != (oauthProvider.GitHubConfig{}) {
		oauthProviders.Add(oauthProvider.NewGitHub(cfg.OAuthGitHub))
	}

	return &adapters{
		oauthProviders: oauthProviders,
		discovery:      discovery,
	}
}
