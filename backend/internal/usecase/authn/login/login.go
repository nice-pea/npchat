package authn

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	Key string

	DB *gorm.DB
}

type Out struct {
	User    model.User    `json:"user"`
	Session model.Session `json:"session"`
}

const sessionLifetime = 24 * time.Hour

func (p Params) Run() (*Out, error) {
	out := &Out{}
	if err := p.DB.Raw(`
		SELECT u.* 
		FROM users AS u
			INNER JOIN credentials AS c 
				ON c.user_id = u.id 
		WHERE c.key = ?`,
		p.Key,
	).Take(&out.User).Error; err != nil {
		return nil, err
	}

	out.Session = model.Session{
		UserID:    out.User.ID,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(sessionLifetime),
	}

	if err := p.DB.Create(&out.Session).Error; err != nil {
		return nil, err
	}

	return out, nil
}
