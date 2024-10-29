package model

import (
	"time"
)

type Chat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	CreatorID uint      `json:"creator_id"`
}
