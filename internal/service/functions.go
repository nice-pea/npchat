package service

import (
	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
)

// getChat возвращает чат либо ошибку ErrChatNotExists
func getChat(repo chatt.Repository, chatID string) (chatt.Chat, error) {
	chatsFilter := domain.ChatsFilter{
		IDs: []string{chatID},
	}
	chats, err := repo.ByChatFilter(chatsFilter)
	if err != nil {
		return chatt.Chat{}, err
	}
	if len(chats) != 1 {
		return chatt.Chat{}, ErrChatNotExists
	}

	return chats[0], nil
}

// getChatByInvitation возвращает чат либо ошибку ErrChatNotExists
func getChatByInvitation(repo chatt.Repository, invitationID string) (chatt.Chat, error) {
	invitationsFilter := domain.InvitationsFilter{
		ID: invitationID,
	}
	chats, err := repo.ByInvitationsFilter(invitationsFilter)
	if err != nil {
		return chatt.Chat{}, err
	}
	if len(chats) != 1 {
		return chatt.Chat{}, ErrChatNotExists
	}

	return chats[0], nil
}
