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

func (s *githubTestSuite) Test_CheckAccess_ValidConfig() {
	provider := NewGithub(GithubConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
	})

	err := provider.CheckAccess()
	s.NoError(err)
}

func (s *githubTestSuite) Test_CheckAccess_EmptyClientID() {
	provider := NewGithub(GithubConfig{
		ClientID:     "",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
	})

	err := provider.CheckAccess()
	s.Error(err)
	s.Contains(err.Error(), "ClientID не может быть пустым")
}

func (s *githubTestSuite) Test_CheckAccess_EmptyClientSecret() {
	provider := NewGithub(GithubConfig{
		ClientID:     "test-client-id",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/callback",
	})

	err := provider.CheckAccess()
	s.Error(err)
	s.Contains(err.Error(), "ClientSecret не может быть пустым")
}

func (s *githubTestSuite) Test_CheckAccess_EmptyRedirectURL() {
	provider := NewGithub(GithubConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "",
	})

	err := provider.CheckAccess()
	s.Error(err)
	s.Contains(err.Error(), "RedirectURL не может быть пустым")
}