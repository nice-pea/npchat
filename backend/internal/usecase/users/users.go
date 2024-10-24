package users

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	DB *gorm.DB
}

func (p Params) Run() (users []model.User, _ error) {
	return users, p.DB.Find(&users).Error
}
