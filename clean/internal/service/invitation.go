package service

import "github.com/saime-0/nice-pea-chat/internal/domain"

type Invitations struct {
	ChatsRepo       domain.ChatsRepository
	MembersRepo     domain.MembersRepository
	InvitationsRepo domain.InvitationsRepository
	History         History
}

type ChatInvitationsInput struct {
	UserID string
	ChatID string
}

// ChatInvitations - возвращает список приглашений данного чата
func (i *Invitations) ChatInvitations(in ChatInvitationsInput) ([]domain.Invitation, error) {
	return nil, nil
}
