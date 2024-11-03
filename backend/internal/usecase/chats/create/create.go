package create

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"

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
		CreatorID: pgtype.Uint32{Uint32: uint32(p.Chat.CreatorID), Valid: true},
	}
	return chat, p.DB.Create(&chat).Error
}
