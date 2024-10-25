package create

import (
	"errors"

	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	Member struct {
		ChatID uint `json:"chat_id"`
		UserID uint `json:"user_id"`
		Pinned bool `json:"pinned"`
	}
	//InvitedByUserID uint `json:"invited_by_user_id"`

	DB *gorm.DB
}

var ErrMemberAlreadyExists = errors.New("member already exists")

func (p Params) Run() (model.Member, error) {
	member := model.Member{
		UserID:   p.Member.UserID,
		ChatID:   p.Member.ChatID,
		IsPinned: p.Member.Pinned,
	}

	var count int64
	if err := p.DB.Model(&member).
		Where("user_id = ?", member.UserID).
		Where("chat_id = ?", member.ChatID).
		Count(&count).Error; err != nil {
		return model.Member{}, err
	} else if count != 0 {
		return model.Member{}, ErrMemberAlreadyExists
	}

	return member, p.DB.Create(&member).Error
}
