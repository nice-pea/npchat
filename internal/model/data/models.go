package data

import (
	"github.com/saime-0/cute-chat-backend/internal/perm"
	"time"
)

type ID string

func (i ID) IsZero() bool {
	return i == ""
}

type User struct {
	ID       ID
	Username string
}

type Credentials struct {
	Login string
}

type Message struct {
	ID       ID
	ChatID   ID
	Date     time.Time
	Text     string
	UserID   ID
	ReplyID  ID
	EditDate time.Time
	DelDate  time.Time
}

type ReplyMessage struct {
	ID       ID
	ChatID   ID
	Date     time.Time
	Text     string
	UserID   ID
	EditDate time.Time
	DelDate  time.Time
}

type Chat struct {
	ID      ID
	OwnerID ID
}

type Member struct {
	ID          ID
	ChatID      ID
	UserID      ID
	Permissions []perm.Kind
}
