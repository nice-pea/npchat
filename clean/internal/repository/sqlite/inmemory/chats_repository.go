package inmemory

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func (m *SQLiteInMemory) NewChatsRepository() (domain.ChatsRepository, error) {
	return &ChatsRepository{
		DB: m.db,
	}, nil
}

type ChatsRepository struct {
	DB *sqlx.DB
}

func (c *ChatsRepository) List(filter domain.ChatsFilter) ([]domain.Chat, error) {
	chats := make([]domain.Chat, 0)
	if err := c.DB.Select(&chats, `
			SELECT * 
			FROM chats 
			WHERE ($1 = "" OR $1 = id)
		`, filter.ID); err != nil {
		return nil, fmt.Errorf("error selecting chats: %w", err)
	}

	return chats, nil
}

func (c *ChatsRepository) Save(chat domain.Chat) error {
	if chat.ID == "" {
		return fmt.Errorf("invalid chat id")
	}
	_, err := c.DB.Exec("INSERT INTO chats(id, name) VALUES (?, ?)", chat.ID, chat.Name)
	if err != nil {
		return fmt.Errorf("error inserting chat: %w", err)
	}

	return nil
}

func (c *ChatsRepository) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("invalid chat id")
	}
	_, err := c.DB.Exec("DELETE FROM chats WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting chat: %w", err)
	}

	return nil
}
