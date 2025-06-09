package oauth_provider

import (
	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

// Mock представляет собой мок для OAuthProvider, используемый для тестирования.
type Mock struct {
	ExchangeFunc         func(code string) (userr.OpenAuthToken, error)              // Функция для обмена кода на токен
	UserFunc             func(token userr.OpenAuthToken) (userr.OpenAuthUser, error) // Функция для получения информации о пользователе
	AuthorizationURLFunc func(state string) string                                   // Функция для генерации URL авторизации
}

// Name возвращает имя провайдера OAuth.
func (m *Mock) Name() string {
	return "mock"
}

// Exchange обменивает код авторизации на токен OAuth.
func (m *Mock) Exchange(code string) (userr.OpenAuthToken, error) {
	// Вызвать мок-функцию, если она определена
	if m.ExchangeFunc != nil {
		return m.ExchangeFunc(code)
	}

	// Паника, если мок-функция не определена
	panic("Exchange not mocked")
}

// User получает информацию о пользователе, используя токен OAuth.
func (m *Mock) User(token userr.OpenAuthToken) (userr.OpenAuthUser, error) {
	// Вызвать мок-функцию, если она определена
	if m.UserFunc != nil {
		return m.UserFunc(token)
	}

	// Паника, если мок-функция не определена
	panic("User not mocked")
}

// AuthorizationURL генерирует URL для авторизации с использованием кода состояния.
func (m *Mock) AuthorizationURL(state string) string {
	// Вызвать мок-функцию, если она определена
	if m.AuthorizationURLFunc != nil {
		return m.AuthorizationURLFunc(state)
	}

	// Паника, если мок-функция не определена
	panic("AuthorizationURL not mocked")
}
