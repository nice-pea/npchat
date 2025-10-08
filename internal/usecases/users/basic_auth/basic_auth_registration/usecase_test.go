package basicAuthRegistration

import (
	"testing"

	"github.com/stretchr/testify/mock"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/common"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	serviceSuite "github.com/nice-pea/npchat/internal/usecases/suite"
)

type testSuite struct {
	serviceSuite.SuiteWithMocks
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) Test_BasicAuthRegistration() {
	usecase := &BasicAuthRegistrationUsecase{
		Repo:         suite.RR.Users,
		SessionsRepo: suite.RR.Sessions,
	}
	mockRepoUsers := suite.RR.Users
	mockRepoSessions := suite.RR.Sessions
	suite.Run("BasicAuthLogin должен быть валидным", func() {
		// Регистрация по логину паролю
		input := In{
			Login:    "",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := usecase.BasicAuthRegistration(input)
		suite.ErrorIs(err, ErrLoginIsRequired)
		suite.Zero(out)
	})

	suite.Run("Password должен быть валидным", func() {
		// Регистрация по логину паролю
		input := In{
			Login:    "login",
			Password: "",
			Name:     "name",
			Nick:     "nick",
		}
		out, err := usecase.BasicAuthRegistration(input)
		suite.ErrorIs(err, ErrPasswordIsRequired)
		suite.Zero(out)
	})

	suite.Run("Name должен быть валидным", func() {
		// Регистрация по логину паролю
		input := In{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "",

			Nick: "nick",
		}
		out, err := usecase.BasicAuthRegistration(input)
		suite.ErrorIs(err, ErrNameIsRequired)
		suite.Zero(out)
	})

	suite.Run("нельзя создать пользователя с существующим логином", func() {
		user := suite.NewRndUserWithBasicAuth()
		input := In{
			Login:    user.BasicAuth.Login,
			Password: common.RndPassword(),
			Name:     "name2",
			Nick:     "nick2",
		}
		mockRepoUsers.EXPECT().InTransaction(mock.Anything).RunAndReturn(func(fn func(userr.Repository) error) error {
			return fn(mockRepoUsers)
		}).Once()
		mockRepoUsers.EXPECT().List(mock.Anything).Return([]userr.User{user}, nil).Once()
		out, err := usecase.BasicAuthRegistration(input)
		suite.ErrorIs(err, ErrLoginIsAlreadyInUse)
		suite.Zero(out)
	})

	suite.Run("после регистрации будет создан пользователь", func() {
		// Регистрация по логину паролю
		input := In{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		// Настройка моков
		mockRepoUsers.EXPECT().InTransaction(mock.Anything).RunAndReturn(func(fn func(userr.Repository) error) error {
			return fn(mockRepoUsers)
		}).Once()
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()
		mockRepoUsers.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		mockRepoSessions.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		out, err := usecase.BasicAuthRegistration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		suite.Equal(input.Name, out.User.Name)
		suite.Equal(input.Nick, out.User.Nick)
	})

	suite.Run("после регистрации будет создана сессия", func() {
		// Регистрация по логину паролю
		input := In{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		var session sessionn.Session
		// Настройка моков
		mockRepoUsers.EXPECT().InTransaction(mock.Anything).RunAndReturn(func(fn func(userr.Repository) error) error {
			return fn(mockRepoUsers)
		}).Once()
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()
		mockRepoUsers.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		mockRepoSessions.EXPECT().Upsert(mock.Anything).Run(func(sessionRepo sessionn.Session) {
			// получаем сессию которая записывается в бд
			session = sessionRepo
		}).Return(nil).Once()
		out, err := usecase.BasicAuthRegistration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Проверить сессию
		suite.EqualSessions(out.Session, session)
		suite.Equal(sessionn.StatusVerified, session.Status)
	})

	suite.Run("после регистрации будет создан метод входа", func() {
		// Регистрация по логину паролю
		input := In{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		// Настройка моков
		mockRepoUsers.EXPECT().InTransaction(mock.Anything).RunAndReturn(func(fn func(userr.Repository) error) error {
			return fn(mockRepoUsers)
		}).Once()
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()
		var user userr.User
		mockRepoUsers.EXPECT().Upsert(mock.Anything).Run(func(userRepo userr.User) {
			user = userRepo
		}).Return(nil).Once()
		mockRepoSessions.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		out, err := usecase.BasicAuthRegistration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		suite.Equal(out.User.ID, user.ID)
		suite.Equal(input.Login, user.BasicAuth.Login)
		suite.Equal(input.Password, user.BasicAuth.Password)
	})
}
