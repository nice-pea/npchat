package chats

import (
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/app/extend"
	"github.com/saime-0/nice-pea-chat/internal/model/rich"
)

type Params struct {
	IDs     []uint
	UserIDs []uint

	DB *gorm.DB
}

type Out struct {
	Chats []rich.Chat `json:"chats,omitempty"`
}

func (p Params) Run() (Out, error) {
	cond := p.DB

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

	// Add creator
	cond = cond.
		Joins(`
		LEFT JOIN users AS creator
			ON creator.id = chats.creator_id`)

	//// Add last_message
	//cond = cond.
	//	Joins(`
	//	LEFT JOIN (
	//		SELECT *, MAX(messages.id) AS last_message_id
	//		FROM messages
	//		LEFT JOIN users AS author
	//			ON author.id = messages.author_id
	//		GROUP BY messages.chat_id
	//	) AS last_message
	//		ON last_message.chat_id = chats.id`)

	var out Out
	if err := cond.Select("*").Find(&out.Chats).Error; err != nil {
		return Out{}, err
	}

	return out, p.extend(out)
}

func (p Params) extend(out Out) error {
	// State
	type State struct {
		DB      *gorm.DB
		ChatIDs []uint
		Chats   map[uint]*rich.Chat

		LastMessagesIDs []uint
		LastMessages    map[uint]*rich.Message

		LastMessagesRepliesIDs []uint
		LastMessagesReplies    map[uint]*rich.MessageReplyTo
	}

	// Collect by IDs
	state := State{
		DB:      p.DB,
		Chats:   make(map[uint]*rich.Chat, len(out.Chats)),
		ChatIDs: make([]uint, len(out.Chats)),
	}
	for i, chat := range out.Chats {
		state.ChatIDs[i] = chat.ID
		state.Chats[chat.ID] = &chat
	}

	// Extending Params
	e := extend.Params[State]{
		State: state,
		Fields: []extend.Field[State]{
			{
				Key: "last_message",
				Fn: func(state State) error {
					lastMsgs := make([]rich.Message, 0, len(state.Chats))
					if err := state.DB.Raw(`
						SELECT *, MAX(messages.id) AS last_message_id FROM messages
						LEFT JOIN users AS author
							ON author.id = messages.author_id
						WHERE messages.chat_id IN (?) 
						GROUP BY messages.chat_id`,
						state.ChatIDs,
					).Scan(&lastMsgs).Error; err != nil {
						return err
					}

					// Collect by IDs
					state.LastMessages = make(map[uint]*rich.Message, len(lastMsgs))
					state.LastMessagesIDs = make([]uint, len(lastMsgs))
					for i, msg := range lastMsgs {
						state.LastMessagesIDs[i] = msg.ID
						state.LastMessages[msg.ID] = &msg
						// Save to chats
						state.Chats[msg.ChatID].LastMessage = &msg
					}

					return nil
				},
			},
			{
				Key:  "last_message_replies",
				Deps: []string{"last_message"},
				Fn: func(state State) error {
					replies := make([]rich.MessageReplyTo, 0, len(state.Chats))
					state.LastMessagesRepliesIDs = make([]uint, 0, len(state.LastMessages))
					for _, msg := range state.LastMessages {
						if msg.ReplyToID != 0 {
							state.LastMessagesRepliesIDs = append(state.LastMessagesRepliesIDs, msg.ReplyToID)
						}
					}
					if err := state.DB.Raw(`
						SELECT * FROM messages
						LEFT JOIN users AS author
							ON author.id = messages.author_id
						WHERE messages.id IN (?)`,
						state.LastMessagesRepliesIDs,
					).Scan(&replies).Error; err != nil {
						return err
					}

					// Collect by IDs
					state.LastMessagesReplies = make(map[uint]*rich.MessageReplyTo, len(replies))
					for _, reply := range replies {
						state.LastMessagesReplies[reply.ID] = &reply
					}
					// Save to lastMessage
					for _, message := range state.LastMessages {
						reply := state.LastMessagesReplies[message.ReplyToID]
						message.ReplyTo = reply
					}

					return nil
				},
			},
		},
	}
	return e.Run()
}
