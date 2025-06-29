package service

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
)

func (suite *testSuite) newRndUserWithSession(sessionStatus string) (out struct {
	User    userr.User
	Session sessionn.Session
}) {
	var err error
	out.User, err = userr.NewUser(gofakeit.Name(), "")
	suite.Require().NoError(err)
	err = suite.rr.users.Upsert(out.User)
	suite.Require().NoError(err)

	out.Session, err = sessionn.NewSession(out.User.ID, gofakeit.ChromeUserAgent(), sessionStatus)
	suite.Require().NoError(err)
	err = suite.rr.sessions.Upsert(out.Session)
	suite.Require().NoError(err)

	return
}

func (suite *testSuite) Test_Sessions_Find() {
	suite.Run("токен должен быть передан", func() {
		for range 10 {
			suite.newRndUserWithSession(sessionn.StatusNew)
		}
		input := SessionsFindIn{
			Token: "",
		}
		sessions, err := suite.ss.sessions.Find(input)
		suite.ErrorIs(err, ErrInvalidToken)
		suite.Empty(sessions)
	})

	suite.Run("вернется пустой список если нет совпадающего токена", func() {
		for range 10 {
			suite.newRndUserWithSession(sessionn.StatusNew)
		}
		input := SessionsFindIn{
			Token: uuid.NewString(),
		}
		sessions, err := suite.ss.sessions.Find(input)
		suite.NoError(err)
		suite.Empty(sessions)
	})

	suite.Run("вернется существующая сессия", func() {
		uws := suite.newRndUserWithSession(sessionn.StatusNew)
		input := SessionsFindIn{
			Token: uws.Session.AccessToken.Token,
		}
		sessions, err := suite.ss.sessions.Find(input)
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.equalSessions(uws.Session, sessions[0])
	})
}
