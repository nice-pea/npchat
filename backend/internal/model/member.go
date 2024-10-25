package model

import "time"

type Member struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	ChatID    uint      `json:"chat_id"`
	IsPinned  bool      `gorm:"type:integer" json:"pinned"`
	CreatedAt time.Time `json:"created_at"`
}
