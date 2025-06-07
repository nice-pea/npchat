package userr

// Repository представляет собой интерфейс для работы с репозиторием пользователей.
type Repository interface {
	List(Filter) ([]User, error)
	Upsert(User) error
}

// Filter представляет собой фильтр для выборки пользователей.
type Filter struct {
	ID                string // ID пользователя для фильтрации
	OAuthUserID       string // Фильтрация по ID пользователя провайдера
	OAuthProvider     string // Фильтрация по провайдеру
	BasicAuthLogin    string // Логин пользователя для фильтрации
	BasicAuthPassword string // Пароль пользователя для фильтрации
}
