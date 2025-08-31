package findSession

import (
	"errors"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
)

type In struct {
	Token string
}

type Out struct {
	Sessions []sessionn.Session
}

var (
	ErrInvalidToken = errors.New("некорректное значение Token")
)

type FindSessionsUsecase struct {
	Repo sessionn.Repository
}

func (s *FindSessionsUsecase) FindSessions(in In) (Out, error) {
	if in.Token == "" {
		return Out{}, ErrInvalidToken
	}

	sessions, err := s.Repo.List(sessionn.Filter{AccessToken: in.Token})
	if err != nil {
		return Out{}, err
	}

	return Out{
		Sessions: sessions,
	}, nil
}
