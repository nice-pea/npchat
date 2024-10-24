package chats

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	IDs     []uint
	UserIDs []uint

	DB *gorm.DB
}

func (p Params) Run() ([]model.Chat, error) {
	var chats []model.Chat
	cond := p.DB
	if p.IDs != nil {
		cond = cond.Where("chats.id IN (?)", p.IDs)
	}
	if p.UserIDs != nil {
		cond = cond.
			Joins("INNER JOIN members ON members.chat_id = chats.id").
			Where("members.user_id IN (?)", p.UserIDs)
	}
	return chats, cond.Find(&chats).Error
}
