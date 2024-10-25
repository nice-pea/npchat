package messages

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/app/optional"
	"github.com/saime-0/nice-pea-chat/internal/model"
)

type Params struct {
	IDs        []uint
	ChatIDs    []uint
	AuthorIDs  []uint
	ReplyToIDs []uint
	Boundary   Boundary
	Limit      optional.Uint

	DB *gorm.DB
}

type Boundary struct {
	AroundID optional.Uint `json:"around_id"`
	BeforeID optional.Uint `json:"before_id"`
	AfterID  optional.Uint `json:"after_id"`
}

func (p Params) Run() ([]model.Message, error) {
	var messages []model.Message
	cond := p.DB
	if p.IDs != nil {
		cond = cond.Where("messages.id IN (?)", p.IDs)
	}
	if p.ChatIDs != nil {
		cond = cond.Where("messages.chat_id IN (?)", p.ChatIDs)
	}
	if p.ChatIDs != nil {
		cond = cond.Where("messages.author_id IN (?)", p.AuthorIDs)
	}
	if p.ReplyToIDs != nil {
		cond = cond.Where("messages.reply_to_id IN (?)", p.ReplyToIDs)
	}

	// Boundary
	if id, ok := p.Boundary.AroundID.Read(); ok {
		messagesTable := p.DB.Raw(`
			WITH below AS (
				SELECT * FROM messages
			   	WHERE id <= @aroundId
			   	ORDER BY id DESC
			   	LIMIT @belowLimit
			),	above AS (
				SELECT * FROM messages
				WHERE id > @aroundId
				LIMIT @aboveLimit
			)
			SELECT * FROM below UNION SELECT * FROM above
		`, map[string]any{
			"aroundId":   id,
			"aboveLimit": p.Limit.Val / 2,
			"belowLimit": p.Limit.Val - (p.Limit.Val / 2),
		})
		cond = cond.Table("(?) messages", messagesTable)
	} else if id, ok = p.Boundary.BeforeID.Read(); ok {
		cond = cond.Where("messages.id < ?", id)
	} else if id, ok = p.Boundary.AfterID.Read(); ok {
		cond = cond.Where("messages.id > ?", id)
	}

	if limit, ok := p.Limit.Read(); ok {
		cond = cond.Limit(int(limit))
	}

	return messages, cond.Find(&messages).Error
}
