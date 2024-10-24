package roles

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	IDs  []uint `json:"ids"`
	Name string `json:"name"`

	DB *gorm.DB
}

func (p Params) Run() (roles []model.Role, _ error) {
	return roles, p.DB.Find(&roles).Error
}
