package service

import (
	"github.com/saime-0/nice-pea-chat/internal/domain/sessionn"
)

// Sessions implements the authentication interface
type Sessions struct {
	Repo sessionn.Repository
}

type SessionsFindIn struct {
	Token string
}

func (s *Sessions) Find(in SessionsFindIn) ([]sessionn.Session, error) {
	if in.Token == "" {
		return nil, ErrInvalidToken
	}

	sessions, err := s.Repo.List(sessionn.Filter{AccessToken: in.Token})
	if err != nil {
		return nil, err
	}

	return sessions, nil
}
