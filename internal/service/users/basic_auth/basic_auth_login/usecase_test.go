package basicAuthLogin

import (
	"testing"

	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	serviceSuite "github.com/nice-pea/npchat/internal/service/suite"
	"github.com/nice-pea/npchat/internal/service/users/oauth"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) Test_BasicAuthLogin() {
	usecase := &BasicAuthLoginUsecase{
		Providers:    oauth.OAuthProviders{},
		Repo:         suite.RR.Users,
		SessionsRepo: suite.RR.Sessions,
	}

	suite.Run("Login должен быть валидным", func() {
		out, err := usecase.BasicAuthLogin(In{
			Login:    " inv ald login",
			Password: "somePassword123!",
		})
		suite.ErrorIs(err, ErrInvalidLogin)
		suite.Zero(out)
	})

	suite.Run("Password должен быть валидным", func() {
		out, err := usecase.BasicAuthLogin(In{
			Login:    "someLogin",
			Password: "invalidpassword",
		})
		suite.ErrorIs(err, ErrInvalidPassword)
		suite.Zero(out)
	})

	suite.Run("неверные данные", func() {
		out, err := usecase.BasicAuthLogin(In{
			Login:    "wrongLogin",
			Password: "wrongPassword123!",
		})
		suite.ErrorIs(err, ErrLoginOrPasswordDoesNotMatch)
		suite.Zero(out)
	})

	suite.Run("вернется Verified сессия", func() {
		// Создаем нового пользователя с AuthnPassword
		user := suite.NewRndUserWithBasicAuth()
		err := suite.RR.Users.Upsert(user)
		if err != nil {
			return
		}
		// Входим сессию с правильными данными
		input := In{
			Login:    user.BasicAuth.Login,
			Password: user.BasicAuth.Password,
		}
		output, err := usecase.BasicAuthLogin(input)
		suite.NoError(err)
		suite.Require().NotZero(output)
		suite.Require().Equal(user.ID, output.Session.UserID)
		suite.Require().Equal(sessionn.StatusVerified, output.Session.Status)

		// Проверяем, что сессия сохранена в репозитории
		sessions, err := suite.RR.Sessions.List(sessionn.Filter{
			AccessToken: output.Session.AccessToken.Token,
		})
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.EqualSessions(output.Session, sessions[0])
	})
}
