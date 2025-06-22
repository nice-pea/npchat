package app

import "github.com/nice-pea/npchat/internal/service"

type services struct {
	chats    *service.Chats
	sessions *service.Sessions
	users    *service.Users
}

func (s *services) Chats() *service.Chats {
	return s.chats
}

func (s *services) Sessions() *service.Sessions {
	return s.sessions
}

func (s *services) Users() *service.Users {
	return s.users
}

func initServices(repos *repositories, adaps *adapters) *services {
	return &services{
		chats: &service.Chats{
			Repo: repos.chats,
		},
		sessions: &service.Sessions{
			Repo: repos.sessions,
		},
		users: &service.Users{
			Providers:    adaps.OAuthProviders(),
			Repo:         repos.users,
			SessionsRepo: repos.sessions,
		},
	}
}
