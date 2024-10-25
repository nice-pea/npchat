package members

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/app/optional"
	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	IDs      []uint
	UserIDs  []uint
	ChatIDs  []uint
	IsPinned optional.Bool

	DB *gorm.DB
}

func (p Params) Run() ([]model.Member, error) {
	cond := p.DB

	if len(p.IDs) > 0 {
		cond = cond.Where("members.id IN (?)", p.IDs)
	}
	if len(p.UserIDs) > 0 {
		cond = cond.Where("members.user_id IN (?)", p.UserIDs)
	}
	if len(p.ChatIDs) > 0 {
		cond = cond.Where("members.chat_id IN (?)", p.ChatIDs)
	}
	if !p.IsPinned.None() {
		cond = cond.Where("members.is_pinned = ?", p.IsPinned.Val())
	}
	var members []model.Member
	return members, cond.Find(&members).Error
}
