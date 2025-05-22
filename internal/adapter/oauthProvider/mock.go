package oauthProvider

import "github.com/saime-0/nice-pea-chat/internal/domain"

const ProviderNameMock = "mock"

// Mock мок для OAuthProvider
type Mock struct {
	ExchangeFunc    func(code string) (domain.OAuthToken, error)
	UserFunc        func(token domain.OAuthToken) (domain.OAuthUser, error)
	AuthCodeURLFunc func(state string) string
}

func (m *Mock) Exchange(code string) (domain.OAuthToken, error) {
	if m.ExchangeFunc != nil {
		return m.ExchangeFunc(code)
	}
	panic("Exchange not mocked")
}

func (m *Mock) User(token domain.OAuthToken) (domain.OAuthUser, error) {
	if m.UserFunc != nil {
		return m.UserFunc(token)
	}
	panic("User not mocked")
}

func (m *Mock) AuthCodeURL(state string) string {
	if m.AuthCodeURLFunc != nil {
		return m.AuthCodeURLFunc(state)
	}
	panic("AuthCodeURL not mocked")
}
