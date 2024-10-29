package chats

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/model/rich"
)

type Params struct {
	IDs     []uint
	UserIDs []uint

	DB *gorm.DB
}

type Out struct {
	Chats []rich.Chat `json:"chats,omitempty"`
}

func (p Params) Run() (Out, error) {
	cond := p.DB

	// Select only with received ids
	if p.IDs != nil {
		cond = cond.Where("chats.id IN (?)", p.IDs)
	}

	// Select only have received members
	if p.UserIDs != nil {
		cond = cond.
			Joins("INNER JOIN members ON members.chat_id = chats.id").
			Where("members.user_id IN (?)", p.UserIDs)
	}

	// Add creator
	cond = cond.
		Joins(`
		LEFT JOIN users AS creator
			ON creator.id = chats.creator_id`)

	var out Out
	if err := cond.Select("*").Find(&out.Chats).Error; err != nil {
		return Out{}, err
	}

	return out, extend(&out, p.DB)
}
