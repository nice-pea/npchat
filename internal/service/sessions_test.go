package service

import (
	"math/rand/v2"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func (suite *servicesTestSuite) newRndUserWithSession(sessionStatus int) (out struct {
	User    domain.User
	Session domain.Session
}) {
	out.User = domain.User{
		ID: uuid.NewString(),
	}
	err := suite.rr.users.Save(out.User)
	suite.Require().NoError(err)

	out.Session = domain.Session{
		ID:     uuid.NewString(),
		UserID: out.User.ID,
		Token:  uuid.NewString(),
		Status: sessionStatus,
	}
	err = suite.rr.sessions.Save(out.Session)
	suite.Require().NoError(err)

	return
}

func (suite *servicesTestSuite) Test_Sessions_Find() {
	suite.Run("токен должен быть передан", func() {
		for range 10 {
			suite.newRndUserWithSession(rand.Int() % 6)
		}
		input := SessionsFindInput{
			Token: "",
		}
		sessions, err := suite.ss.sessions.Find(input)
		suite.ErrorIs(err, ErrInvalidToken)
		suite.Empty(sessions)
	})
	suite.Run("вернется пустой список если нет совпадающего токена", func() {
		for range 10 {
			suite.newRndUserWithSession(rand.Int() % 6)
		}
		input := SessionsFindInput{
			Token: uuid.NewString(),
		}
		sessions, err := suite.ss.sessions.Find(input)
		suite.NoError(err)
		suite.Empty(sessions)
	})
	suite.Run("вернется существующая сессия", func() {
		uws := suite.newRndUserWithSession(rand.Int() % 6)
		input := SessionsFindInput{
			Token: uws.Session.Token,
		}
		sessions, err := suite.ss.sessions.Find(input)
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.Equal(uws.Session, sessions[0])
	})
	//suite.Run("необходимый заголовок должен находится в карте", func() {
	//	input := ByHttpHeader{
	//		Headers: headers{
	//			"other-header": "1234567890",
	//		},
	//	}
	//	user, err := suite.sessionsService.ByHttpHeader(input)
	//	suite.ErrorIs(err, ErrInvalidAuthorizationToken)
	//	suite.Zero(user)
	//})
	//
	//suite.Run("значение заголовка должно иметь определенный формат", func() {
	//	input := ByHttpHeader{
	//		Headers: headers{
	//			sessionAuthHeader: "IOCUgnxcausUIASd 1234567890",
	//		},
	//	}
	//	user, err := suite.sessionsService.ByHttpHeader(input)
	//	suite.ErrorIs(err, ErrInvalidAuthorizationToken)
	//	suite.Zero(user)
	//})
	//
	//suite.Run("токен должен соответствовать пользователю", func() {
	//	input := ByHttpHeader{
	//		Headers: headers{
	//			sessionAuthHeader: sessionAuthScheme + " 1234567890",
	//		},
	//	}
	//	user, err := suite.sessionsService.ByHttpHeader(input)
	//	suite.ErrorIs(err, ErrInvalidAuthorizationToken)
	//	suite.Zero(user)
	//})

	//suite.Run("токен должен соответствовать пользователю", func() {
	//	uws := suite.newRndUserWithSession(SessionStatusVerified)
	//	input := SessionsFindInput{
	//		Token: uws.Session.Token,
	//	}
	//	sessions, err := suite.ss.sessions.Find(input)
	//	suite.NoError(err)
	//	suite.Require().Len(sessions, 1)
	//	suite.Equal(uws.Session, sessions[0])
	//
	//	input := SessionGet{
	//		Headers: headers{
	//			sessionAuthHeader: sessionAuthScheme + " 1234567890",
	//		},
	//	}
	//	user, err := suite.sessionsService.ByHttpHeader(input)
	//	suite.ErrorIs(err, ErrInvalidAuthorizationToken)
	//	suite.Zero(user)
	//})
	//suite.sessionsService.Repository

	//OAuth{
	//	StateUnq
	//	UserCode
	//	Token
	//	ExpiresAt
	//	RefreshToken
	//}

	//session {
	//	ID
	//	userID
	//	Token
	//	Status [New|Pending|Verified|Expired|Revoked|Failed]
	//}
	//
	//[Credentials|OAuthGoogle|OAuthVK]SessionVerification{ // ..._SV{}
	//	SessionID
	//	...
	//}
	//suite.Run("вернется существующая сессия", func() {
	//	uws := suite.newRndUserWithSession(SessionStatusVerified)
	//	input := SessionsFindInput{
	//		Token: uws.Session.Token,
	//	}
	//	sessions, err := suite.ss.sessions.Find(input)
	//	suite.NoError(err)
	//	suite.Require().Len(sessions, 1)
	//	suite.Equal(uws.Session, sessions[0])
	//})
}
