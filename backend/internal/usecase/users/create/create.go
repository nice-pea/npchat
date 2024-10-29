package create

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	User struct {
		Username string `json:"username"`
		Key      string `json:"key"`
	}

	DB *gorm.DB
}

func (p Params) Run() (model.User, error) {
	user := model.User{
		Username: p.User.Username,
	}
	credentials := model.Credentials{
		UserID: 0,
		Key:    p.User.Key,
	}

	return user, p.DB.Transaction(func(tx *gorm.DB) error {
		if err := p.DB.Create(&user).Error; err != nil {
			return err
		}
		credentials.UserID = user.ID

		return p.DB.Create(&credentials).Error
	})
}
