package domain

type Chat struct {
	ID   string
	Name string
}

type ChatsRepository interface {
	List(filter ChatsFilter) ([]Chat, error)
	Save(chat Chat) error
	Delete(id string) error
}

type ChatsFilter struct {
	ID      string
	UserIDs []string
}
