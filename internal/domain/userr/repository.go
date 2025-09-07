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
	OauthUserID       string    // Фильтрация по ID пользователя провайдера
	OauthProvider     string    // Фильтрация по провайдеру
	BasicAuthLogin    string    // Логин пользователя для фильтрации
	BasicAuthPassword string    // Пароль пользователя для фильтрации
}

// Find возвращает пользователя либо ошибку ErrUserNotExists
func Find(repo Repository, filter Filter) (User, error) {
	users, err := repo.List(filter)
	if err != nil {
		return User{}, ErrUserNotExists
	}
	if len(users) != 1 {
		return User{}, ErrUserNotExists
	}

	return users[0], nil
}
