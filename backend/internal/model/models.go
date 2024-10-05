package model

import (
	"github.com/saime-0/cute-chat-backend/internal/perm"
	"time"
)

/*
	todo
	- Модели данных
	- функции валидации
	- интерфейсы policy
*/

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

// todo: rm and use fet calls usecases from controller, and build one api model
// Deprecated: delete
type ChatOfMember struct {
	  ID ID
	  Name string
	  LastMsg Message?
	  LastReadMsg ID // optional
	  UnreadCount int
	  Permissions []perm.Kind
}
