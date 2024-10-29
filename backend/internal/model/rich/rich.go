package rich

import (
	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Chat struct {
	model.Chat  `gorm:"embedded"`
	Creator     *model.User `gorm:"embedded" json:"creator,omitempty"`
	LastMessage *Message    `json:"last_message,omitempty"`
}

type Message struct {
	model.Message `gorm:"embedded"`
	Author        *model.User     `gorm:"embedded" json:"author,omitempty"`
	ReplyTo       *MessageReplyTo `gorm:"embedded" json:"reply_to,omitempty"`
}

type MessageReplyTo struct {
	model.Message `gorm:"embedded"`
	Author        model.User `gorm:"embedded" json:"author"`
}
