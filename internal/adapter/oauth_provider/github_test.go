package oauthProvider

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type githubTestSuite struct {
	suite.Suite
}

func Test_GithubTestSuite(t *testing.T) {
	suite.Run(t, new(githubTestSuite))
}

func (s *githubTestSuite) Test_NewGithub() {
	s.Run("валидная конфигурация", func() {
		provider, err := NewGithub(GithubConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
		})

		s.NoError(err)
		s.NotNil(provider)
	})

	s.Run("пустой ClientID", func() {
		provider, err := NewGithub(GithubConfig{
			ClientID:     "",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
		})

		s.Error(err)
		s.Nil(provider)
		s.Contains(err.Error(), "ClientID не может быть пустым")
	})

	s.Run("пустой ClientSecret", func() {
		provider, err := NewGithub(GithubConfig{
			ClientID:     "test-client-id",
			ClientSecret: "",
			RedirectURL:  "http://localhost:8080/callback",
		})

		s.Error(err)
		s.Nil(provider)
		s.Contains(err.Error(), "ClientSecret не может быть пустым")
	})

	s.Run("пустой RedirectURL", func() {
		provider, err := NewGithub(GithubConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "",
		})

		s.Error(err)
		s.Nil(provider)
		s.Contains(err.Error(), "RedirectURL не может быть пустым")
	})
}