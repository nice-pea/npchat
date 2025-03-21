package base

import "github.com/saime-0/nice-pea-chat/internal/domain"

var _ domain.ChatsRepository = (*ChatsRepository)(nil)

type ChatsRepository struct {
}

func (c *ChatsRepository) List(filter domain.ChatsFilter) ([]domain.Chat, error) {
	return []domain.Chat{}, nil
}

func (c *ChatsRepository) Save(chat domain.Chat) error {
	return nil
}

func (c *ChatsRepository) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}
