package model

import (
	"time"

	"github.com/saime-0/nice-pea-chat/internal/app/null"
)

type Chat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	CreatorID null.Uint `json:"creator_id"`
}
