package domain

import "errors"

type AuthnPassword struct {
	UserID   string
	Login    string
	Password string
}

func (c AuthnPassword) ValidateLogin() error {
	if c.Login == "" {
		return errors.New("login is empty")
	}

	return nil
}

type AuthnPasswordRepository interface {
	Save(AuthnPassword) error
	List(filter AuthnPasswordFilter) ([]AuthnPassword, error)
}

type AuthnPasswordFilter struct {
	UserID   string
	Login    string
	Password string
}
