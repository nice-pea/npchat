package sessionn

type Repository interface {
	List(Filter) ([]Session, error)
	Upsert(Session) error
	InTransaction(func(txRepo Repository) error) error
}

// Filter представляет собой фильтр по сессиям.
type Filter struct {
	AccessToken string // Фильтрация по токену сессии
}
