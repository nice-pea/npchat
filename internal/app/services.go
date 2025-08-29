package app

import (
	"github.com/nice-pea/npchat/internal/service"
	sessionsFind "github.com/nice-pea/npchat/internal/service/sessions_find"
)

type services struct {
	chats *service.Chats
	*sessionsFind.SessionsFindUsecase
	users *service.Users
}

func (s *services) Chats() *service.Chats {
	return s.chats
}

func (s *services) Users() *service.Users {
	return s.users
}

func initServices(rr *repositories, aa *adapters) *services {
	return &services{
		chats: &service.Chats{
			Repo: rr.chats,
		},
		SessionsFindUsecase: &sessionsFind.SessionsFindUsecase{
			Repo: rr.sessions,
		},
		users: &service.Users{
			Providers:    aa.OAuthProviders(),
			Repo:         rr.users,
			SessionsRepo: rr.sessions,
		},
	}
}
