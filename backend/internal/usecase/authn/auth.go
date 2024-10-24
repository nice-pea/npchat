package authn

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	Token string `json:"token"`

	DB *gorm.DB
}

type Out struct {
	User    model.User    `gorm:"embedded" json:"user"`
	Session model.Session `gorm:"embedded" json:"session"`
}

func (p Params) Run() (*Out, error) {
	out := &Out{}
	if err := p.DB.Raw(`
		SELECT u.*, s.*
		FROM users AS u
			INNER JOIN sessions AS s 
				ON s.user_id = u.id 
		WHERE s.token = ?`,
		p.Token,
	).Take(&out).Error; err != nil {
		return nil, err
	}

	return out, nil
}
