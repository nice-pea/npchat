package rich

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/saime-0/nice-pea-chat/internal/model"
)

// ChatRow - структура для использования с pgx.RowToStructByPos
type ChatRow struct {
	ChatID        uint
	ChatName      string
	ChatCreatedAt time.Time
	ChatCreatorID pgtype.Uint32

	// Chat.Creator [Optional]
	CreatorID        pgtype.Uint32
	CreatorUsername  pgtype.Text
	CreatorCreatedAt pgtype.Timestamp

	// Chat.LastMessage [Optional]
	LastMsgID        pgtype.Uint32
	LastMsgChatID    pgtype.Uint32
	LastMsgText      pgtype.Text
	LastMsgAuthorID  pgtype.Uint32
	LastMsgReplyToID pgtype.Uint32
	LastMsgEditedAt  pgtype.Timestamp
	LastMsgRemovedAt pgtype.Timestamp
	LastMsgCreatedAt pgtype.Timestamp

	// Chat.LastMessage.Author [Optional]
	LastMsgAuthorID1       pgtype.Uint32
	LastMsgAuthorUsername  pgtype.Text
	LastMsgAuthorCreatedAt pgtype.Timestamp

	// Chat.LastMessage.Reply [Optional]
	LastMsgReplyID        pgtype.Uint32
	LastMsgReplyChatID    pgtype.Uint32
	LastMsgReplyText      pgtype.Text
	LastMsgReplyAuthorID  pgtype.Uint32
	LastMsgReplyReplyToID pgtype.Uint32
	LastMsgReplyEditedAt  pgtype.Timestamp
	LastMsgReplyRemovedAt pgtype.Timestamp
	LastMsgReplyCreatedAt pgtype.Timestamp

	// Chat.LastMessage.Reply.Author [Optional]
	LastMsgReplyAuthorID1       pgtype.Uint32
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

	if row.ChatCreatorID.Valid {
		chat.Creator = &model.User{
			ID:        uint(row.CreatorID.Uint32),
			Username:  row.CreatorUsername.String,
			CreatedAt: row.CreatorCreatedAt.Time,
		}
	}

	if row.LastMsgID.Valid {
		chat.LastMessage = &Message{
			Message: model.Message{
				ID:        uint(row.LastMsgID.Uint32),
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
		if row.LastMsgAuthorID.Valid {
			chat.LastMessage.Author = &model.User{
				ID:        uint(row.LastMsgAuthorID.Uint32),
				Username:  row.LastMsgAuthorUsername.String,
				CreatedAt: row.LastMsgAuthorCreatedAt.Time,
			}
		}
		if row.LastMsgReplyID.Valid {
			chat.LastMessage.ReplyTo = &Reply{
				Message: model.Message{
					ID:        uint(row.LastMsgReplyID.Uint32),
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
			if row.LastMsgReplyAuthorID.Valid {
				chat.LastMessage.ReplyTo.Author = &model.User{
					ID:        uint(row.LastMsgReplyAuthorID.Uint32),
					Username:  row.LastMsgReplyAuthorUsername.String,
					CreatedAt: row.LastMsgReplyAuthorCreatedAt.Time,
				}
			}
		}
	}

	return chat
}
