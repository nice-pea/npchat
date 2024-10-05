package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/saime-0/cute-chat-backend/internal/model"
	"github.com/saime-0/cute-chat-backend/internal/repository"
)

var _ repository.Common = (*CommonRepository)(nil)

type CommonRepository struct {
	conn *pgx.Conn
}

func NewCommonRepository(conn *pgx.Conn) *CommonRepository {
	return &CommonRepository{conn: conn}
}

func (p *CommonRepository) UserCreate(ctx context.Context, user model.User) error {
	//TODO implement me
	panic("implement me")
}

func (p *CommonRepository) UserUpdate(ctx context.Context, user model.User) error {
	//TODO implement me
	panic("implement me")
}

func (p *CommonRepository) ChatCreate(ctx context.Context, chat model.Chat) error {
	//TODO implement me
	panic("implement me")
}

func (p *CommonRepository) ChatUpdate(ctx context.Context, chat model.Chat) error {
	//TODO implement me
	panic("implement me")
}

func (p *CommonRepository) MemberCreate(ctx context.Context, member model.Member) error {
	//TODO implement me
	panic("implement me")
}

func (p *CommonRepository) MemberUpdate(ctx context.Context, member model.Member) error {
	//TODO implement me
	panic("implement me")
}

func (p *CommonRepository) MemberDelete(ctx context.Context, id model.ID) error {
	//TODO implement me
	panic("implement me")
}

func (p *CommonRepository) UserByID(ctx context.Context, id model.ID) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (p *CommonRepository) UserByUserCreds(ctx context.Context, credentials model.Credentials) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (p *CommonRepository) ChatByID(ctx context.Context, id model.ID) (*model.Chat, error) {
	//TODO implement me
	panic("implement me")
}

func (p *CommonRepository) MemberByChatID(ctx context.Context, id model.ID) ([]model.Member, error) {
	//TODO implement me
	panic("implement me")
}

func (p *CommonRepository) MemberByUserID(ctx context.Context, id model.ID) ([]model.Member, error) {
	//TODO implement me
	panic("implement me")
}
