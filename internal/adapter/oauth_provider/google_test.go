package oauthProvider

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type googleTestSuite struct {
	suite.Suite
}

func Test_GoogleTestSuite(t *testing.T) {
	suite.Run(t, new(googleTestSuite))
}

func (s *googleTestSuite) Test_NewGoogle() {
	s.Run("валидная конфигурация", func() {
		provider, err := NewGoogle(GoogleConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
		})

		s.NoError(err)
		s.NotNil(provider)
	})

	s.Run("пустой ClientID", func() {
		provider, err := NewGoogle(GoogleConfig{
			ClientID:     "",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
		})

		s.Error(err)
		s.Nil(provider)
		s.Contains(err.Error(), "ClientID не может быть пустым")
	})

	s.Run("пустой ClientSecret", func() {
		provider, err := NewGoogle(GoogleConfig{
			ClientID:     "test-client-id",
			ClientSecret: "",
			RedirectURL:  "http://localhost:8080/callback",
		})

		s.Error(err)
		s.Nil(provider)
		s.Contains(err.Error(), "ClientSecret не может быть пустым")
	})

	s.Run("пустой RedirectURL", func() {
		provider, err := NewGoogle(GoogleConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "",
		})

		s.Error(err)
		s.Nil(provider)
		s.Contains(err.Error(), "RedirectURL не может быть пустым")
	})
}