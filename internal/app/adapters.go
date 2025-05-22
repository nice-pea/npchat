package app

import (
	"os"

	"github.com/saime-0/nice-pea-chat/internal/adapter"
	"github.com/saime-0/nice-pea-chat/internal/adapter/oauthProvider"
	"github.com/saime-0/nice-pea-chat/internal/service"
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

func initAdapters() *adapters {
	var discovery = &adapter.ServiceDiscoveryBase{
		Debug: true,
	}

	return &adapters{
		oauthProviders: service.OAuthProviders{
			oauthProvider.ProviderNameGoogle: &oauthProvider.Google{
				ClientID:     os.Getenv("GOOGLE_KEY"),
				ClientSecret: os.Getenv("GOOGLE_SECRET"),
				RedirectURL:  discovery.NpcApiPubUrl() + "/oauth/google/registration/callback",
			},
			oauthProvider.ProviderNameGitHub: &oauthProvider.GitHub{
				ClientID:     os.Getenv("GITHUB_KEY"),
				ClientSecret: os.Getenv("GITHUB_SECRET"),
				RedirectURL:  discovery.NpcApiPubUrl() + "/oauth/github/registration/callback",
			},
		},
		discovery: discovery,
	}
}
