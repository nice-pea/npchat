package chats

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/app/null"
	"github.com/saime-0/nice-pea-chat/internal/app/optional"
	"github.com/saime-0/nice-pea-chat/internal/model/rich"
)

type Params struct {
	IDs                  []uint
	UserIDs              []uint
	UnreadCounterForUser optional.Uint
	UpdateBefore         null.Time
	DB                   *gorm.DB
}

type Out struct {
	Chats []rich.Chat `json:"chats,omitempty"`
}

func (p Params) Run() (Out, error) {
	cond := p.DB.Raw(`
		WITH c1 AS (
			SELECT
				c.id,
				COALESCE(MAX(m.id), 0) AS last_msg_id
			FROM chats c
				LEFT JOIN messages m
					ON c.id = m.chat_id
			GROUP BY c.id
			ORDER BY last_msg_id DESC
		)
		SELECT chats.*, creator.*, last_msg.*, msg_author.*, reply_msg.*, reply_author.*
		FROM c1
			INNER JOIN chats 
					ON c1.id = chats.id
			LEFT JOIN users AS creator 
					ON chats.creator_id = creator.id
			LEFT JOIN messages last_msg
				   ON c1.last_msg_id = last_msg.id
			LEFT JOIN users AS msg_author
					ON last_msg.author_id = msg_author.id
			LEFT JOIN messages AS reply_msg
					ON last_msg.reply_to_id = reply_msg.id
			LEFT JOIN users AS reply_author
					ON reply_msg.author_id = reply_author.id
		`,
	)

	// Select only with received ids
	if p.IDs != nil {
		cond = cond.Where("chats.id IN (?)", p.IDs)
	}

	// Select only have received members
	if p.UserIDs != nil {
		cond = cond.
			Joins("INNER JOIN members ON members.chat_id = chats.id").
			Where("members.user_id IN (?)", p.UserIDs)
	}

	var out Out
	if err := cond.Scan(&out.Chats).Error; err != nil {
		return Out{}, err
	}

	return out, extend(&out, p)
}
