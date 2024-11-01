package rich

import (
	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Chat struct {
	model.Chat          `gorm:"embedded;embeddedPrefix:chats_"`
	Creator             *model.User `gorm:"embedded;embeddedPrefix:creator_" json:"creator,omitempty"`
	LastMessage         *Message    `gorm:"embedded;embeddedPrefix:last_msg_" json:"last_message,omitempty"`
	UnreadMessagesCount int         `json:"unread_messages_count,omitempty"`
}

type Message struct {
	model.Message `gorm:"embedded"`
	Author        *model.User     `gorm:"embedded;embeddedPrefix:author_" json:"author,omitempty"`
	ReplyTo       *MessageReplyTo `gorm:"embedded;embeddedPrefix:reply_" json:"reply_to,omitempty"`
}

type MessagesMap map[uint]*Message

func MsgsMap(msgs []Message) MessagesMap {
	msgsMap := make(MessagesMap, len(msgs))
	for _, msg := range msgs {
		msgsMap[msg.ID] = &msg
	}

	return msgsMap
}

type MessageReplyTo struct {
	model.Message `gorm:"embedded"`
	Author        *model.User `gorm:"embedded;embeddedPrefix:author_" json:"author"`
}
