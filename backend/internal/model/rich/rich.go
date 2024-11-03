package rich

import (
	"github.com/saime-0/nice-pea-chat/internal/model"
)

//func ChatScanDst() (*Chat, []any) {
//	ch := &Chat{
//		Chat:    model.Chat{},
//		Creator: &model.User{},
//		LastMessage: &Message{
//			Message: model.Message{},
//			Author:  &model.User{},
//			ReplyTo: &Reply{
//				Message: model.Message{},
//				Author:  &model.User{},
//			},
//		},
//		UnreadMessagesCount: 0,
//	}
//	return ch, []any{
//		&ch.ID,
//		&ch.Name,
//		&ch.CreatedAt,
//		&ch.CreatorID,
//
//		&ch.Creator.ID,
//		&ch.Creator.Username,
//		&ch.Creator.CreatedAt,
//
//		&ch.LastMessage.ID,
//		&ch.LastMessage.ChatID,
//		&ch.LastMessage.Text,
//		&ch.LastMessage.AuthorID,
//		&ch.LastMessage.ReplyToID,
//		&ch.LastMessage.EditedAt,
//		&ch.LastMessage.RemovedAt,
//		&ch.LastMessage.CreatedAt,
//
//		&ch.LastMessage.Author.ID,
//		&ch.LastMessage.Author.Username,
//		&ch.LastMessage.Author.CreatedAt,
//
//		&ch.LastMessage.ReplyTo.ID,
//		&ch.LastMessage.ReplyTo.ChatID,
//		&ch.LastMessage.ReplyTo.Text,
//		&ch.LastMessage.ReplyTo.AuthorID,
//		&ch.LastMessage.ReplyTo.ReplyToID,
//		&ch.LastMessage.ReplyTo.EditedAt,
//		&ch.LastMessage.ReplyTo.RemovedAt,
//		&ch.LastMessage.ReplyTo.CreatedAt,
//
//		&ch.LastMessage.ReplyTo.Author.ID,
//		&ch.LastMessage.ReplyTo.Author.Username,
//		&ch.LastMessage.ReplyTo.Author.CreatedAt,
//	}
//}

// rich model

type Chat struct {
	model.Chat
	Creator             *model.User `json:"creator,omitempty"`
	LastMessage         *Message    `json:"last_message,omitempty"`
	UnreadMessagesCount int         `json:"unread_messages_count,omitempty"`
}
type Message struct {
	model.Message
	Author  *model.User `json:"author,omitempty"`
	ReplyTo *Reply      `json:"reply_to,omitempty"`
}

type Reply struct {
	model.Message
	Author *model.User `json:"author,omitempty"`
}
