package create

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/app/null"
	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	Message struct {
		ChatID    uint      `json:"chat_id"`
		Text      string    `json:"text"`
		AuthorID  null.Uint `json:"author_id"`
		ReplyToID null.Uint `json:"reply_to_id"`
	}

	DB *gorm.DB
}

func (p Params) Run() (model.Message, error) {
	chat := model.Message{
		ChatID:    p.Message.ChatID,
		Text:      p.Message.Text,
		AuthorID:  p.Message.AuthorID,
		ReplyToID: p.Message.ReplyToID,
	}
	return chat, p.DB.Create(&chat).Error
}
