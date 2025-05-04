package domain

import "errors"

type LoginCredentials struct {
	UserID   string
	Login    string
	Password string
}

func (c LoginCredentials) ValidateLogin() error {
	if c.Login == "" {
		return errors.New("login is empty")
	}

	return nil
}

type LoginCredentialsRepository interface {
	Save(LoginCredentials) error
	List(filter LoginCredentialsFilter) ([]LoginCredentials, error)
}

type LoginCredentialsFilter struct {
	UserID   string
	Login    string
	Password string
}
