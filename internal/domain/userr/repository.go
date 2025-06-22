package userr

import "github.com/google/uuid"

// Repository представляет собой интерфейс для работы с репозиторием пользователей.
type Repository interface {
	List(Filter) ([]User, error)
	Upsert(User) error
	InTransaction(func(txRepo Repository) error) error
}

// Filter представляет собой фильтр для выборки пользователей.
type Filter struct {
	ID                uuid.UUID // ID пользователя для фильтрации
	OAuthUserID       string    // Фильтрация по ID пользователя провайдера
	OAuthProvider     string    // Фильтрация по провайдеру
	BasicAuthLogin    string    // Логин пользователя для фильтрации
	BasicAuthPassword string    // Пароль пользователя для фильтрации
}
