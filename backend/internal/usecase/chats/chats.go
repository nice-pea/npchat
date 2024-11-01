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
	c1SubQuery := p.DB.Table("chats").
		Select("chats.id, COALESCE(MAX(messages.id), 0) AS last_msg_id").
		Joins("LEFT JOIN messages ON chats.id = messages.chat_id").
		Group("chats.id").
		Order("last_msg_id DESC")

	cond := p.DB.Table("(?) as c1", c1SubQuery).
		Select(`
			chats.id AS chats_id,
			chats.name AS chats_name,
			chats.created_at AS chats_created_at,
			chats.creator_id AS chats_creator_id,

			creator.id AS creator_id, 
			creator.username AS creator_username, 
			creator.created_at AS creator_created_at, 

			last_msg.id AS last_msg_id,
			last_msg.chat_id AS last_msg_chat_id,
			last_msg.text AS last_msg_text,
			last_msg.author_id AS last_msg_author_id,
			last_msg.reply_to_id AS last_msg_reply_to_id,
			last_msg.edited_at AS last_msg_edited_at,
			last_msg.removed_at AS last_msg_removed_at,
			last_msg.created_at AS last_msg_created_at,

			last_msg_author.id AS last_msg_author_id, 
			last_msg_author.username AS last_msg_author_username, 
			last_msg_author.created_at AS last_msg_author_created_at,

			last_msg_reply.id AS last_msg_reply_id,
			last_msg_reply.chat_id AS last_msg_reply_chat_id,
			last_msg_reply.text AS last_msg_reply_text,
			last_msg_reply.author_id AS last_msg_reply_author_id,
			last_msg_reply.reply_to_id AS last_msg_reply_reply_to_id,
			last_msg_reply.edited_at AS last_msg_reply_edited_at,
			last_msg_reply.removed_at AS last_msg_reply_removed_at,
			last_msg_reply.created_at AS last_msg_reply_created_at,

			last_msg_reply_author.id AS last_msg_reply_author_id, 
			last_msg_reply_author.username AS last_msg_reply_author_username, 
			last_msg_reply_author.created_at AS last_msg_reply_author_created_at
		`).
		Joins(`
			INNER JOIN chats
					ON c1.id = chats.id
			LEFT JOIN users AS creator
					ON chats.creator_id = creator.id
			LEFT JOIN messages last_msg
				   ON c1.last_msg_id = last_msg.id
			LEFT JOIN users AS last_msg_author
					ON last_msg.author_id = last_msg_author.id
			LEFT JOIN messages AS last_msg_reply
					ON last_msg.reply_to_id = last_msg_reply.id
			LEFT JOIN users AS last_msg_reply_author
					ON last_msg_reply.author_id = last_msg_reply_author.id
		`)

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
