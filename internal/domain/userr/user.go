package userr

type User struct {
	ID   string // ID пользователя
	Name string // Имя пользователя
	Nick string // Ник пользователя

	BasicAuth     BasicAuth      // Данные для аутентификации по логину и паролю
	OpenAuthLinks []OpenAuthLink // Связи для аутентификации по OAuth
}
