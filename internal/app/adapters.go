package app

import (
	"os"

	"golang.org/x/oauth2"

	googleOAuth "golang.org/x/oauth2/google"

	"github.com/saime-0/nice-pea-chat/internal/adapter"
)

type adapters struct {
	oauthGoogle adapter.OAuthGoogle
	discovery   adapter.ServiceDiscovery
}

func (a *adapters) Discovery() adapter.ServiceDiscovery {
	return a.discovery
}

func (a *adapters) OAuthGoogle() adapter.OAuthGoogle {
	return a.oauthGoogle
}

func initAdapters() *adapters {
	var discovery = &adapter.ServiceDiscoveryBase{}
	authGoogleBase := &adapter.OAuthGoogleBase{
		Config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_KEY"),
			ClientSecret: os.Getenv("GOOGLE_SECRET"),
			Endpoint:     googleOAuth.Endpoint,
			RedirectURL:  discovery.NpcApiPubUrl() + "/oauth/google/registration/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
	}
	return &adapters{
		oauthGoogle: authGoogleBase,
		discovery:   discovery,
	}
}
