package repository

import (
	"context"
	"github.com/saime-0/cute-chat-backend/internal/model/data"
)

//}
//type DomainRead interface {

type Domain interface {
	UserCreate(context.Context, data.User) error
	UserUpdate(context.Context, data.User) error

	ChatCreate(context.Context, data.Chat) error
	ChatUpdate(context.Context, data.Chat) error

	MemberCreate(context.Context, data.Member) error
	MemberUpdate(context.Context, data.Member) error
	MemberDelete(context.Context, data.ID) error

	UserByID(context.Context, data.ID) (*data.User, error)
	UserByUserCreds(context.Context, data.Credentials) (*data.User, error)

	ChatByID(context.Context, data.ID) (*data.Chat, error)

	MemberByChatID(context.Context, data.ID) ([]data.Member, error)
	MemberByUserID(context.Context, data.ID) ([]data.Member, error)
}

type Data interface {
	UserCreate(context.Context, data.User) error
	UserUpdate(context.Context, data.User) error

	ChatCreate(context.Context, data.Chat) error
	ChatUpdate(context.Context, data.Chat) error

	MemberCreate(context.Context, data.Member) error
	MemberUpdate(context.Context, data.Member) error
	MemberDelete(context.Context, data.ID) error

	UserByID(context.Context, data.ID) (*data.User, error)
	UserByUserCreds(context.Context, data.Credentials) (*data.User, error)

	ChatByID(context.Context, data.ID) (*data.Chat, error)

	MemberByChatID(context.Context, data.ID) ([]data.Member, error)
	MemberByUserID(context.Context, data.ID) ([]data.Member, error)
}
