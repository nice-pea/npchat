package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/nullism/bqb"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type chat struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	ChiefUserID string `db:"chief_user_id"`
}

func chatToDomain(repoChat chat) domain.Chat {
	return domain.Chat{
		ID:          repoChat.ID,
		Name:        repoChat.Name,
		ChiefUserID: repoChat.ChiefUserID,
	}
}

func chatsToDomain(repoChats []chat) []domain.Chat {
	domainChats := make([]domain.Chat, len(repoChats))
	for i, repoChat := range repoChats {
		domainChats[i] = chatToDomain(repoChat)
	}

	return domainChats
}

func (m *RepositoryFactory) NewChatsRepository() (domain.ChatsRepository, error) {
	return &ChatsRepository{
		DB: m.db,
	}, nil
}

type ChatsRepository struct {
	DB *sqlx.DB
}

func (c *ChatsRepository) List(filter domain.ChatsFilter) ([]domain.Chat, error) {
	// Построить запрос используя bqb
	where := bqb.Optional("WHERE")
	if len(filter.IDs) > 0 {
		where.And("id IN (?)", filter.IDs)
	}
	sql, args, err := bqb.New("SELECT * FROM chats ?", where).ToSql()
	if err != nil {
		return nil, err
	}
	// Выполнить запрос используя sqlx
	chats := make([]chat, 0)
	if err = c.DB.Select(&chats, sql, args...); err != nil {
		return nil, fmt.Errorf("error selecting chats: %w", err)
	}

	return chatsToDomain(chats), nil
}

func (c *ChatsRepository) Save(chat domain.Chat) error {
	if chat.ID == "" {
		return fmt.Errorf("invalid chat id")
	}
	_, err := c.DB.Exec(`
		INSERT OR REPLACE INTO chats(id, name, chief_user_id)
		VALUES (?, ?, ?)`,
		chat.ID, chat.Name, chat.ChiefUserID)
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
