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

func (s *googleTestSuite) Test_CheckAccess_ValidConfig() {
	provider := NewGoogle(GoogleConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
	})

	err := provider.CheckAccess()
	s.NoError(err)
}

func (s *googleTestSuite) Test_CheckAccess_EmptyClientID() {
	provider := NewGoogle(GoogleConfig{
		ClientID:     "",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
	})

	err := provider.CheckAccess()
	s.Error(err)
	s.Contains(err.Error(), "ClientID не может быть пустым")
}

func (s *googleTestSuite) Test_CheckAccess_EmptyClientSecret() {
	provider := NewGoogle(GoogleConfig{
		ClientID:     "test-client-id",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/callback",
	})

	err := provider.CheckAccess()
	s.Error(err)
	s.Contains(err.Error(), "ClientSecret не может быть пустым")
}

func (s *googleTestSuite) Test_CheckAccess_EmptyRedirectURL() {
	provider := NewGoogle(GoogleConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "",
	})

	err := provider.CheckAccess()
	s.Error(err)
	s.Contains(err.Error(), "RedirectURL не может быть пустым")
}