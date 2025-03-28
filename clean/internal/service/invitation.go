package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type Invitations struct {
	ChatsRepo       domain.ChatsRepository
	MembersRepo     domain.MembersRepository
	InvitationsRepo domain.InvitationsRepository
	History         History
}

var (
	ErrChatInvitationsInputUserIDValidate = errors.New("некорректный UserID")
	ErrChatInvitationsInputChatIDValidate = errors.New("некорректный ChatID")
)

type ChatInvitationsInput struct {
	UserID string
	ChatID string
}

func (in ChatInvitationsInput) validate() error {
	if err := uuid.Validate(in.ChatID); err != nil {
		return errors.Join(err, ErrChatInvitationsInputChatIDValidate)
	}
	if err := uuid.Validate(in.UserID); err != nil {
		return errors.Join(err, ErrChatInvitationsInputUserIDValidate)
	}
	return nil
}

// ChatInvitations - возвращает список приглашений данного чата
func (i *Invitations) ChatInvitations(in ChatInvitationsInput) ([]domain.Invitation, error) {
	return nil, nil
}
