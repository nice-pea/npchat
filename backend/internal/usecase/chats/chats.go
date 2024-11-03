package chats

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/nullism/bqb"

	"github.com/saime-0/nice-pea-chat/internal/app/null"
	"github.com/saime-0/nice-pea-chat/internal/app/optional"
	"github.com/saime-0/nice-pea-chat/internal/model/rich"
)

type Params struct {
	IDs                  []uint
	UserIDs              []uint
	UnreadCounterForUser optional.Uint
	UpdateBefore         null.Time
	Conn                 *pgx.Conn
}

type Out struct {
	Chats []rich.Chat `json:"chats,omitempty"`
}

func (p Params) Run() (Out, error) {
	sel := bqb.New(`
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
	`)
	where := bqb.Optional("WHERE")

	// Select only with received ids
	if p.IDs != nil {
		where.And("chats.id = ANY(?)", p.IDs)
	}

	// Select only have received members
	if p.UserIDs != nil {
		sel.Space(`INNER JOIN members ON members.chat_id = chats.id`)
		where.And("members.user_id = ANY(?)", p.UserIDs)
	}

	query, args, err := bqb.New("? ?", sel, where).ToPgsql()
	if err != nil {
		return Out{}, fmt.Errorf("to_pgsql: %w", err)
	}

	// Run query
	rows, err := p.Conn.Query(context.Background(), query, args...)
	if err != nil {
		return Out{}, err
	}

	// Scan into DTO
	var chatRows rich.ChatRows
	chatRows, err = pgx.CollectRows[rich.ChatRow](rows, pgx.RowToStructByPos)
	if err != nil {
		return Out{}, err
	}

	// Result from DTO
	out := Out{
		Chats: chatRows.Rich(),
	}

	return out, extend(&out, p)
}
