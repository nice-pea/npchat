package model

import (
	"time"

	"github.com/saime-0/nice-pea-chat/internal/app/null"
)

type Chat struct {
	ID        uint      `gorm:"primaryKey" json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	CreatorID null.Uint `json:"creator_id" db:"creator_id"`
}
