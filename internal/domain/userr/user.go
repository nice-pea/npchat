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

// Equal проверяет пользователей на равенство
func (u *User) Equal(u2 User) bool {
	if u.ID != u2.ID {
		return false
	}
	if u.Name != u2.Name {
		return false
	}
	if u.Nick != u2.Nick {
		return false
	}
	if len(u2.OpenAuthUsers) != len(u.OpenAuthUsers) {
		return false
	}

	for i, openAuthUser := range u2.OpenAuthUsers {
		if openAuthUser.ID != u.OpenAuthUsers[i].ID {
			return false
		}
		if openAuthUser.Provider != u.OpenAuthUsers[i].Provider {
			return false
		}
		if openAuthUser.Email != u.OpenAuthUsers[i].Email {
			return false
		}
		if openAuthUser.Name != u.OpenAuthUsers[i].Name {
			return false
		}
		if openAuthUser.Picture != u.OpenAuthUsers[i].Picture {
			return false
		}

		token2 := openAuthUser.Token
		token1 := u.OpenAuthUsers[i].Token
		if token1.AccessToken != token2.AccessToken {
			return false
		}
		if token1.TokenType != token2.TokenType {
			return false
		}
		if token1.RefreshToken != token2.RefreshToken {
			return false
		}
		if !token2.Expiry.Equal(token1.Expiry) {
			return false
		}
	}

	return true
}
