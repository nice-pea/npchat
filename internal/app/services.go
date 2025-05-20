package app

import "github.com/saime-0/nice-pea-chat/internal/service"

type services struct {
	chats         *service.Chats
	invitations   *service.Invitations
	members       *service.Members
	sessions      *service.Sessions
	authnPassword *service.AuthnPassword
	oauth         *service.OAuth
}

func (s *services) OAuth() *service.OAuth {
	return s.oauth
}

func (s *services) Chats() *service.Chats {
	return s.chats
}

func (s *services) Invitations() *service.Invitations {
	return s.invitations
}

func (s *services) Members() *service.Members {
	return s.members
}

func (s *services) Sessions() *service.Sessions {
	return s.sessions
}

func (s *services) AuthnPassword() *service.AuthnPassword {
	return s.authnPassword
}

func initServices(repos *repositories, adaps *adapters) *services {
	return &services{
		chats: &service.Chats{
			ChatsRepo:   repos.chats,
			MembersRepo: repos.members,
		},
		invitations: &service.Invitations{
			ChatsRepo:       repos.chats,
			MembersRepo:     repos.members,
			InvitationsRepo: repos.invitations,
			UsersRepo:       repos.users,
		},
		members: &service.Members{
			MembersRepo: repos.members,
			ChatsRepo:   repos.chats,
		},
		sessions: &service.Sessions{
			SessionsRepo: repos.sessions,
		},
		authnPassword: &service.AuthnPassword{
			AuthnPasswordRepo: repos.authnPassword,
			SessionsRepo:      repos.sessions,
		},
		oauth: &service.OAuth{
			Google: adaps.oauthGoogle,
		},
	}
}
