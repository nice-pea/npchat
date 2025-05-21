package app

import (
	"os"

	"github.com/saime-0/nice-pea-chat/internal/adapter"
)

type adapters struct {
	oauthProviders adapter.OAuthProviders
	discovery      adapter.ServiceDiscovery
}

func (a *adapters) Discovery() adapter.ServiceDiscovery {
	return a.discovery
}

func (a *adapters) OAuthProviders() adapter.OAuthProviders {
	return a.oauthProviders
}

func initAdapters() *adapters {
	var discovery = &adapter.ServiceDiscoveryBase{
		Debug: true,
	}

	return &adapters{
		oauthProviders: adapter.OAuthProviders{
			adapter.OAuthProviderGoogle: &adapter.OAuthGoogle{
				ClientID:     os.Getenv("GOOGLE_KEY"),
				ClientSecret: os.Getenv("GOOGLE_SECRET"),
				RedirectURL:  discovery.NpcApiPubUrl() + "/oauth/google/registration/callback",
			},
			adapter.OAuthProviderGithub: &adapter.OAuthGitHub{
				ClientID:     os.Getenv("GITHUB_KEY"),
				ClientSecret: os.Getenv("GITHUB_SECRET"),
				RedirectURL:  discovery.NpcApiPubUrl() + "/oauth/github/registration/callback",
			},
		},
		discovery: discovery,
	}
}
