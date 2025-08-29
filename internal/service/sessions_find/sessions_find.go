package sessionsfind

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

type Interface interface {
	SessionsFind(In) (Out, error)
}

type Impl struct {
	Repo sessionn.Repository
}

func (s *Impl) SessionsFind(in In) (Out, error) {
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

// type Func = func(In) (Out, error)

// func New(repo sessionn.Repository) Func {
// 	return func(in In) (Out, error) {
// 		if in.Token == "" {
// 			return Out{}, ErrInvalidToken
// 		}

// 		sessions, err := repo.List(sessionn.Filter{AccessToken: in.Token})
// 		if err != nil {
// 			return Out{}, err
// 		}

// 		return Out{
// 			Sessions: sessions,
// 		}, nil
// 	}
// }
