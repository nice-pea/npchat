package repository

import (
	"context"
	"github.com/saime-0/cute-chat-backend/internal/model"
)

type Common interface {
	UserCreate(context.Context, model.User) error
	UserUpdate(context.Context, model.User) error

	ChatCreate(context.Context, model.Chat) error
	ChatUpdate(context.Context, model.Chat) error

	MemberCreate(context.Context, model.Member) error
	MemberUpdate(context.Context, model.Member) error
	MemberDelete(context.Context, model.ID) error

	UserByID(context.Context, model.ID) (*model.User, error)
	UserByUserCreds(context.Context, model.Credentials) (*model.User, error)

	ChatByID(context.Context, model.ID) (*model.Chat, error)

	MemberByChatID(context.Context, model.ID) ([]model.Member, error)
	MemberByUserID(context.Context, model.ID) ([]model.Member, error)
}
