package basicAuthLogin

import (
	"testing"

	"github.com/stretchr/testify/mock"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	serviceSuite "github.com/nice-pea/npchat/internal/usecases/suite"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) Test_BasicAuthLogin() {
	usecase := &BasicAuthLoginUsecase{
		Repo:         suite.RR.Users,
		SessionsRepo: suite.RR.Sessions,
	}
	mockRepoUsers := suite.RR.Users
	mockRepoSessions := suite.RR.Sessions

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
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()
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
		// Входим сессию с правильными данными
		input := In{
			Login:    user.BasicAuth.Login,
			Password: user.BasicAuth.Password,
		}
		mockRepoUsers.EXPECT().List(mock.Anything).Return([]userr.User{user}, nil).Once()
		mockRepoSessions.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		output, err := usecase.BasicAuthLogin(input)
		suite.NoError(err)
		suite.Require().NotZero(output)
		suite.Require().Equal(user.ID, output.Session.UserID)
		suite.Require().Equal(sessionn.StatusVerified, output.Session.Status)
	})
}
