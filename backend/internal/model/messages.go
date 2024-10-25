package model

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ChatID    uint      `json:"chat_id"`
	Text      string    `json:"text"`
	AuthorID  uint      `json:"author_id"`
	ReplyToID uint      `json:"reply_to_id"`
	EditedAt  time.Time `json:"edited_at"`
	RemovedAt time.Time `json:"removed_at"`
	CreatedAt time.Time `json:"created_at"`
}
