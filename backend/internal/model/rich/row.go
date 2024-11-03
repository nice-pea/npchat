package rich

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/saime-0/nice-pea-chat/internal/app/null"
	"github.com/saime-0/nice-pea-chat/internal/model"
)

// ChatRow - структура для использования с pgx.RowToStructByPos
type ChatRow struct {
	ChatID        uint
	ChatName      string
	ChatCreatedAt time.Time
	ChatCreatorID null.Uint

	// Chat.Creator [Optional]
	CreatorID        null.Uint
	CreatorUsername  pgtype.Text
	CreatorCreatedAt pgtype.Timestamp

	// Chat.LastMessage [Optional]
	LastMsgID        null.Uint
	LastMsgChatID    null.Uint
	LastMsgText      pgtype.Text
	LastMsgAuthorID  null.Uint
	LastMsgReplyToID null.Uint
	LastMsgEditedAt  pgtype.Timestamp
	LastMsgRemovedAt pgtype.Timestamp
	LastMsgCreatedAt pgtype.Timestamp

	// Chat.LastMessage.Author [Optional]
	LastMsgAuthorID1       null.Uint
	LastMsgAuthorUsername  pgtype.Text
	LastMsgAuthorCreatedAt pgtype.Timestamp

	// Chat.LastMessage.Reply [Optional]
	LastMsgReplyID        null.Uint
	LastMsgReplyChatID    null.Uint
	LastMsgReplyText      pgtype.Text
	LastMsgReplyAuthorID  null.Uint
	LastMsgReplyReplyToID null.Uint
	LastMsgReplyEditedAt  pgtype.Timestamp
	LastMsgReplyRemovedAt pgtype.Timestamp
	LastMsgReplyCreatedAt pgtype.Timestamp

	// Chat.LastMessage.Reply.Author [Optional]
	LastMsgReplyAuthorID1       null.Uint
	LastMsgReplyAuthorUsername  pgtype.Text
	LastMsgReplyAuthorCreatedAt pgtype.Timestamp
}

type ChatRows []ChatRow

func (rows ChatRows) Rich() []Chat {
	chats := make([]Chat, len(rows))
	for i, row := range rows {
		chats[i] = row.Rich()
	}
	return chats
}

// Rich - функция для преобразования ChatRow в структуру Chat
func (row ChatRow) Rich() Chat {
	chat := Chat{
		Chat: model.Chat{
			ID:        row.ChatID,
			Name:      row.ChatName,
			CreatedAt: row.ChatCreatedAt,
			CreatorID: row.CreatorID,
		},
		Creator:             nil,
		LastMessage:         nil,
		UnreadMessagesCount: 0,
	}

	if row.ChatCreatorID.Valid() {
		chat.Creator = &model.User{
			ID:        row.CreatorID.Val(),
			Username:  row.CreatorUsername.String,
			CreatedAt: row.CreatorCreatedAt.Time,
		}
	}

	if row.LastMsgID.Valid() {
		chat.LastMessage = &Message{
			Message: model.Message{
				ID:        row.LastMsgID.Val(),
				ChatID:    row.ChatID,
				Text:      row.LastMsgText.String,
				AuthorID:  row.LastMsgAuthorID,
				ReplyToID: row.LastMsgReplyReplyToID,
				EditedAt:  row.LastMsgEditedAt,
				RemovedAt: row.LastMsgRemovedAt,
				CreatedAt: row.LastMsgCreatedAt.Time,
			},
			Author:  nil,
			ReplyTo: nil,
		}
		if row.LastMsgAuthorID.Valid() {
			chat.LastMessage.Author = &model.User{
				ID:        row.LastMsgAuthorID.Val(),
				Username:  row.LastMsgAuthorUsername.String,
				CreatedAt: row.LastMsgAuthorCreatedAt.Time,
			}
		}
		if row.LastMsgReplyID.Valid() {
			chat.LastMessage.ReplyTo = &Reply{
				Message: model.Message{
					ID:        row.LastMsgReplyID.Val(),
					ChatID:    row.ChatID,
					Text:      row.LastMsgReplyText.String,
					AuthorID:  row.LastMsgReplyAuthorID,
					ReplyToID: row.LastMsgReplyReplyToID,
					EditedAt:  row.LastMsgReplyEditedAt,
					RemovedAt: row.LastMsgReplyRemovedAt,
					CreatedAt: row.LastMsgReplyCreatedAt.Time,
				},
				Author: nil,
			}
			if row.LastMsgReplyAuthorID.Valid() {
				chat.LastMessage.ReplyTo.Author = &model.User{
					ID:        row.LastMsgReplyAuthorID.Val(),
					Username:  row.LastMsgReplyAuthorUsername.String,
					CreatedAt: row.LastMsgReplyAuthorCreatedAt.Time,
				}
			}
		}
	}

	return chat
}
