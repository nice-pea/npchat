package service

import (
	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// Sessions implements the authentication interface
type Sessions struct {
	SessionsRepo domain.SessionsRepository
}

type SessionsFindInput struct {
	Token string
}

func (s *Sessions) Find(in SessionsFindInput) ([]domain.Session, error) {
	if in.Token == "" {
		return nil, ErrInvalidToken
	}

	sessions, err := s.SessionsRepo.List(domain.SessionFilter{Token: in.Token})
	if err != nil {
		return nil, err
	}

	return sessions, nil
}
