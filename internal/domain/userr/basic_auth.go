package userr

// BasicAuth представляет собой метод аутентификации по логину и паролю.
type BasicAuth struct {
	Login    string // Логин пользователя
	Password string // Пароль пользователя
}

// NewBasicAuth создает новый метод аутентификации по логину и паролю.
func NewBasicAuth(login string, password string) (BasicAuth, error) {
	if err := ValidateBasicAuthLogin(login); err != nil {
		return BasicAuth{}, err
	}
	if err := ValidateBasicAuthPassword(password); err != nil {
		return BasicAuth{}, err
	}

	return BasicAuth{
		Login:    login,
		Password: password,
	}, nil
}
