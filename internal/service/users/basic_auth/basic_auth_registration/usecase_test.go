package basicAuthRegistration

import (
	"testing"

	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/common"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	serviceSuite "github.com/nice-pea/npchat/internal/service/suite"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) Test_BasicAuthRegistration() {
	usecase := &BasicAuthRegistrationUsecase{
		Repo:         suite.RR.Users,
		SessionsRepo: suite.RR.Sessions,
	}

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
		// Регистрация по логину паролю
		input := In{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := usecase.BasicAuthRegistration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Регистрация второй раз с существующим логином
		input = In{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name2",
			Nick:     "nick2",
		}
		out, err = usecase.BasicAuthRegistration(input)
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
		out, err := usecase.BasicAuthRegistration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		suite.Equal(input.Name, out.User.Name)
		suite.Equal(input.Nick, out.User.Nick)

		// Прочитать пользователя из репозитория
		users, err := suite.RR.Users.List(userr.Filter{})
		suite.Require().NoError(err)
		suite.Require().Len(users, 1)
		suite.Equal(out.User, users[0])
	})

	suite.Run("после регистрации будет создана сессия", func() {
		// Регистрация по логину паролю
		input := In{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := usecase.BasicAuthRegistration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Проверить сессию
		sessions, err := suite.RR.Sessions.List(sessionn.Filter{})
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.EqualSessions(out.Session, sessions[0])
		suite.Equal(sessionn.StatusVerified, sessions[0].Status)
	})

	suite.Run("после регистрации будет создан метод входа", func() {
		// Регистрация по логину паролю
		input := In{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := usecase.BasicAuthRegistration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Прочитать метод входа из репозитория
		users, err := suite.RR.Users.List(userr.Filter{
			BasicAuthLogin: input.Login,
		})
		suite.Require().NoError(err)
		suite.Require().Len(users, 1)
		suite.Equal(out.User.ID, users[0].ID)
		suite.Equal(input.Login, users[0].BasicAuth.Login)
		suite.Equal(input.Password, users[0].BasicAuth.Password)
	})
}
