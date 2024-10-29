package model

import (
	"time"

	"github.com/saime-0/nice-pea-chat/internal/app/null"
)

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ChatID    uint      `json:"chat_id"`
	Text      string    `json:"text"`
	AuthorID  null.Uint `gorm:"default:null" json:"author_id,omitempty"`
	ReplyToID null.Uint `gorm:"default:null" json:"reply_to_id,omitempty"`
	EditedAt  null.Time `gorm:"default:null" json:"edited_at,omitempty"`
	RemovedAt null.Time `gorm:"default:null" json:"removed_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
