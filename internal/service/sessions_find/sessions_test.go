package sessionsFind

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	serviceSuite "github.com/nice-pea/npchat/internal/service/service_suite"
)

func (suite *testSuite) newRndUserWithSession(sessionStatus string) (out struct {
	User    userr.User
	Session sessionn.Session
}) {
	var err error
	out.User, err = userr.NewUser(gofakeit.Name(), "")
	suite.Require().NoError(err)
	err = suite.RR.Users.Upsert(out.User)
	suite.Require().NoError(err)

	out.Session, err = sessionn.NewSession(out.User.ID, gofakeit.ChromeUserAgent(), sessionStatus)
	suite.Require().NoError(err)
	err = suite.RR.Sessions.Upsert(out.Session)
	suite.Require().NoError(err)

	return
}

type testSuite struct {
	serviceSuite.Suite
	*SessionsFindUsecase
}

func (suite *testSuite) SetupTest() {
	suite.Suite.SetupTest()
	suite.SessionsFindUsecase = &SessionsFindUsecase{Repo: suite.RR.Sessions}
}

func (suite *testSuite) Test_Sessions_Find() {
	suite.Run("токен должен быть передан", func() {
		for range 10 {
			suite.newRndUserWithSession(sessionn.StatusNew)
		}
		input := In{
			Token: "",
		}

		out, err := suite.SessionsFind(input)
		suite.ErrorIs(err, ErrInvalidToken)
		suite.Zero(out)
	})

	suite.Run("вернется пустой список если нет совпадающего токена", func() {
		for range 10 {
			suite.newRndUserWithSession(sessionn.StatusNew)
		}
		input := In{
			Token: uuid.NewString(),
		}
		out, err := suite.SessionsFind(input)
		suite.NoError(err)
		suite.Zero(out)
	})

	suite.Run("вернется существующая сессия", func() {
		uws := suite.newRndUserWithSession(sessionn.StatusNew)
		input := In{
			Token: uws.Session.AccessToken.Token,
		}
		out, err := suite.SessionsFind(input)
		suite.NoError(err)
		suite.Require().Len(out.Sessions, 1)
		suite.EqualSessions(uws.Session, out.Sessions[0])
	})
}
