package chatt

// Repository представляет собой интерфейс для работы с репозиторием чатов.
type Repository interface {
	List(Filter) ([]Chat, error)
	Upsert(Chat) error
}

// Filter представляет собой фильтр для выборки чатов.
type Filter struct {
	ID                    string // Фильтрация по ID чата
	InvitationID          string // Фильтрация по ID приглашений в чате
	InvitationRecipientID string // Фильтрация по ID получателей приглашения в чат
	ParticipantID         string // Фильтрация по ID участников в чате
}

// Find возвращает чат либо ошибку ErrChatNotExists
func Find(repo Repository, filter Filter) (Chat, error) {
	chats, err := repo.List(filter)
	if err != nil {
		return Chat{}, err
	}
	if len(chats) != 1 {
		return Chat{}, ErrChatNotExists
	}

	return chats[0], nil
}
