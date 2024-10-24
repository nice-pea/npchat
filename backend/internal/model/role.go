package model

import (
	"github.com/saime-0/nice-pea-chat/internal/app/database/typ"
)

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"-" json:"name"`
	ChatID      uint      `json:"chat_id"`
	Permissions typ.Uints `gorm:"type:text" json:"permissions,omitempty"`
}
