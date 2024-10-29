package chats

import (
	"gorm.io/gorm"

	extend1 "github.com/saime-0/nice-pea-chat/internal/app/extend"
	"github.com/saime-0/nice-pea-chat/internal/model/rich"
)

type chatsExt struct {
	chatIDs []uint
	chats   map[uint]*rich.Chat
	db      *gorm.DB
}

func (e *chatsExt) lastMessages() error {
	lastMsgs := make([]rich.Message, 0, len(e.chats))
	if err := e.db.Raw(`
		SELECT DISTINCT ON (messages.chat_id) *
		FROM messages
			LEFT JOIN users AS author
				ON author.id = messages.author_id
			LEFT JOIN messages AS reply
				ON reply.id = messages.reply_to_id
			LEFT JOIN users AS reply_author
				ON reply_author.id = reply.author_id
		WHERE messages.chat_id IN (?) 
		ORDER BY messages.chat_id, messages.id DESC`,
		e.chatIDs,
	).Scan(&lastMsgs).Error; err != nil {
		return err
	}

	// Save into chatsExt
	for _, msg := range rich.MsgsMap(lastMsgs) {
		e.chats[msg.ChatID].LastMessage = msg
	}

	return nil
}

func extend(out *Out, db *gorm.DB) error {
	ext := &chatsExt{
		db:      db,
		chats:   make(map[uint]*rich.Chat, len(out.Chats)),
		chatIDs: make([]uint, len(out.Chats)),
	}

	// Fill required fields for extending
	for i, chat := range out.Chats {
		ext.chatIDs[i] = chat.ID
		ext.chats[chat.ID] = &chat
	}

	// Extending Params
	extending := extend1.Params{
		Fields: []extend1.Field{
			{
				Key: "last_message",
				Fn:  ext.lastMessages,
			},
		},
	}
	if err := extending.Run(); err != nil {
		return err
	}

	out.Chats = make([]rich.Chat, 0, len(ext.chats))
	for _, chat := range ext.chats {
		out.Chats = append(out.Chats, *chat)
	}

	return nil
}
