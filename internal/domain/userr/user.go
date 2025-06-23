package userr

import (
	"errors"

	"github.com/google/uuid"
)

// User представляет собой агрегат пользователя.
type User struct {
	ID   uuid.UUID // ID пользователя
	Name string    // Имя пользователя
	Nick string    // Ник пользователя

	BasicAuth     BasicAuth      // Данные для аутентификации по логину и паролю
	OpenAuthUsers []OpenAuthUser // Связи для аутентификации по OAuth
}

// NewUser создает нового пользователя с указанным именем и ником.
func NewUser(name string, nick string) (User, error) {
	if err := ValidateUserName(name); err != nil {
		return User{}, err
	}
	if nick != "" {
		if err := ValidateUserNick(nick); err != nil {
			return User{}, err
		}
	}

	return User{
		ID:            uuid.New(),
		Name:          name,
		Nick:          nick,
		BasicAuth:     BasicAuth{},
		OpenAuthUsers: []OpenAuthUser{},
	}, nil
}

// AddOpenAuthUser добавляет нового пользователя в список связей для аутентификации по OAuth.
func (u *User) AddOpenAuthUser(newOpenAuthUser OpenAuthUser) error {
	for _, ou := range u.OpenAuthUsers {
		if ou.Provider == newOpenAuthUser.Provider {
			return errors.New("пользователь уже связан с этим провайдером")
		}
	}

	u.OpenAuthUsers = append(u.OpenAuthUsers, newOpenAuthUser)

	return nil
}

// isBasicAuthSet проверяет, установлен ли метод аутентификации по логину и паролю.
func (u *User) isBasicAuthSet() bool {
	return u.BasicAuth != (BasicAuth{})
}

// AddBasicAuth добавляет пользователю метод аутентификации по логину и паролю.
func (u *User) AddBasicAuth(auth BasicAuth) error {
	if u.isBasicAuthSet() {
		return errors.New("метод аутентификации по логину и паролю уже установлен")
	}

	u.BasicAuth = auth

	return nil
}
