package model

import "time"

type Member struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	ChatID    uint      `gorm:"primaryKey" json:"chat_id"`
	CreatedAt time.Time `json:"created_at"`
}
