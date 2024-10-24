package users

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	IDs []uint

	DB *gorm.DB
}

func (p Params) Run() (users []model.User, _ error) {
	cond := p.DB
	if len(p.IDs) > 0 {
		cond = cond.Where("id IN (?)", p.IDs)
	}

	return users, cond.Find(&users).Error
}
