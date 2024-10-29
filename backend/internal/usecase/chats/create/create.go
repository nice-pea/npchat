package create

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/app/null"
	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	Chat struct {
		Name      string `json:"name"`
		CreatorID uint   `json:"creator_id"`
	}

	DB *gorm.DB
}

func (p Params) Run() (model.Chat, error) {
	chat := model.Chat{
		Name:      p.Chat.Name,
		CreatorID: null.Uint{V: p.Chat.CreatorID},
	}
	return chat, p.DB.Create(&chat).Error
}
