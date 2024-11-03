package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Message struct {
	ID        uint             `json:"id"`
	ChatID    uint             `json:"chat_id"`
	Text      string           `json:"text"`
	AuthorID  pgtype.Uint32    `json:"author_id,omitempty"`
	ReplyToID pgtype.Uint32    `json:"reply_to_id,omitempty"`
	EditedAt  pgtype.Timestamp `json:"edited_at,omitempty"`
	RemovedAt pgtype.Timestamp `json:"removed_at,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}
