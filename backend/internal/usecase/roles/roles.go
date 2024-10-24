package roles

import (
	"gorm.io/gorm"

	. "github.com/saime-0/nice-pea-chat/internal/model/role"
)

type Params struct {
	IDs  []uint `json:"ids"`
	Name string `json:"name"`

	DB *gorm.DB
}

func (p Params) Run() (roles []Role, _ error) {
	return roles, p.DB.Find(&roles).Error
}
