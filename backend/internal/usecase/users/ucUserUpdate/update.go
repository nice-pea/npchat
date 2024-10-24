package ucUserUpdate

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	User struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}

	DB *gorm.DB
}

func (p Params) Run() (model.User, error) {
	user := model.User{
		ID:       p.User.ID,
		Username: p.User.Username,
	}

	return user, p.DB.Updates(&user).Take(&user).Error
}
