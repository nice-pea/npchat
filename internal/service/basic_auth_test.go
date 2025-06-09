package service

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/saime-0/nice-pea-chat/internal/domain/common"
	"github.com/saime-0/nice-pea-chat/internal/domain/sessionn"
	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

func (suite *servicesTestSuite) newRndUserWithBasicAuth() userr.User {
	user, err := userr.NewUser(gofakeit.Name(), gofakeit.Noun())
	suite.Require().NoError(err)
	ba, err := userr.NewBasicAuth(gofakeit.Noun(), common.RndPassword())
	suite.Require().NoError(err)
	err = user.AddBasicAuth(ba)
	suite.Require().NoError(err)
	return user
}

func (suite *servicesTestSuite) Test_BasicAuthLogin() {
	suite.Run("Login должен быть валидным", func() {
		out, err := suite.ss.users.BasicAuthLogin(BasicAuthLoginIn{
			Login:    " inv ald login",
			Password: "somePassword123!",
		})
		suite.ErrorIs(err, ErrInvalidLogin)
		suite.Zero(out)
	})

	suite.Run("Password должен быть валидным", func() {
		out, err := suite.ss.users.BasicAuthLogin(BasicAuthLoginIn{
			Login:    "someLogin",
			Password: "invalidpassword",
		})
		suite.ErrorIs(err, ErrInvalidPassword)
		suite.Zero(out)
	})

	suite.Run("неверные данные", func() {
		out, err := suite.ss.users.BasicAuthLogin(BasicAuthLoginIn{
			Login:    "wrongLogin",
			Password: "wrongPassword123!",
		})
		suite.ErrorIs(err, ErrLoginOrPasswordDoesNotMatch)
		suite.Zero(out)
	})

	suite.Run("вернется Verified сессия", func() {
		// Создаем нового пользователя с AuthnPassword
		user := suite.newRndUserWithBasicAuth()
		// Входим сессию с правильными данными
		input := BasicAuthLoginIn{
			Login:    user.BasicAuth.Login,
			Password: user.BasicAuth.Password,
		}
		output, err := suite.ss.users.BasicAuthLogin(input)
		suite.NoError(err)
		suite.Require().NotZero(output)
		suite.Require().Equal(user.ID, output.Session.UserID)
		suite.Require().Equal(sessionn.StatusVerified, output.Session.Status)

		// Проверяем, что сессия сохранена в репозитории
		sessions, err := suite.rr.sessions.List(sessionn.Filter{
			AccessToken: output.Session.AccessToken.Token,
		})
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.Equal(output.Session, sessions[0])
	})
}

func (suite *servicesTestSuite) Test_AuthnPassword_Registration() {
	suite.Run("BasicAuthLogin должен быть валидным", func() {
		// Регистрация по логину паролю
		input := BasicAuthRegistrationIn{
			Login:    "",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.users.BasicAuthRegistration(input)
		suite.ErrorIs(err, ErrInvalidLogin)
		suite.Zero(out)
	})

	suite.Run("Password должен быть валидным", func() {
		// Регистрация по логину паролю
		input := BasicAuthRegistrationIn{
			Login:    "login",
			Password: "",
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.users.BasicAuthRegistration(input)
		suite.ErrorIs(err, ErrInvalidPassword)
		suite.Zero(out)
	})

	suite.Run("Name должен быть валидным", func() {
		// Регистрация по логину паролю
		input := BasicAuthRegistrationIn{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "",

			Nick: "nick",
		}
		out, err := suite.ss.users.BasicAuthRegistration(input)
		suite.ErrorIs(err, ErrInvalidName)
		suite.Zero(out)
	})

	suite.Run("нельзя создать пользователя с существующим логином", func() {
		// Регистрация по логину паролю
		input := BasicAuthRegistrationIn{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.users.BasicAuthRegistration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Регистрация второй раз с существующим логином
		input = BasicAuthRegistrationIn{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name2",
			Nick:     "nick2",
		}
		out, err = suite.ss.users.BasicAuthRegistration(input)
		suite.ErrorIs(err, ErrLoginIsAlreadyInUse)
		suite.Zero(out)
	})

	suite.Run("после регистрации будет создан пользователь", func() {
		// Регистрация по логину паролю
		input := BasicAuthRegistrationIn{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.users.BasicAuthRegistration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		suite.Equal(input.Name, out.User.Name)
		suite.Equal(input.Nick, out.User.Nick)

		// Прочитать пользователя из репозитория
		users, err := suite.rr.users.List(userr.Filter{})
		suite.Require().NoError(err)
		suite.Require().Len(users, 1)
		suite.Equal(out.User, users[0])
	})

	suite.Run("после регистрации будет создана сессия", func() {
		// Регистрация по логину паролю
		input := BasicAuthRegistrationIn{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.users.BasicAuthRegistration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Проверить сессию
		sessions, err := suite.rr.sessions.List(sessionn.Filter{})
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.Equal(out.Session, sessions[0])
		suite.Equal(sessionn.StatusVerified, sessions[0].Status)
	})

	suite.Run("после регистрации будет создан метод входа", func() {
		// Регистрация по логину паролю
		input := BasicAuthRegistrationIn{
			Login:    "login",
			Password: common.RndPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.users.BasicAuthRegistration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Прочитать метод входа из репозитория
		users, err := suite.rr.users.List(userr.Filter{
			BasicAuthLogin: input.Login,
		})
		suite.Require().NoError(err)
		suite.Require().Len(users, 1)
		suite.Equal(out.User.ID, users[0].ID)
		suite.Equal(input.Login, users[0].BasicAuth.Login)
		suite.Equal(input.Password, users[0].BasicAuth.Password)
	})
}
