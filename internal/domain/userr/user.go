package userr

// User представляет собой агрегат пользователя.
type User struct {
	ID   string // ID пользователя
	Name string // Имя пользователя
	Nick string // Ник пользователя

	BasicAuth     BasicAuth      // Данные для аутентификации по логину и паролю
	OpenAuthLinks []OpenAuthLink // Связи для аутентификации по OAuth
}

// NewUser создает нового пользователя с указанным именем и ником.
func NewUser(name string, nick string) (User, error) {
	if err := ValidateUserName(name); err != nil {
		return User{}, err
	}
	if err := ValidateUserNick(nick); err != nil {
		return User{}, err
	}

	return User{
		Name: name,
		Nick: nick,
	}, nil
}
