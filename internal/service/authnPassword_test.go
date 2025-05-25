package service

import (
	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func (suite *servicesTestSuite) newRndUserWithAuthnPassword() domain.AuthnPassword {
	ap := domain.AuthnPassword{
		UserID:   uuid.NewString(),
		Login:    randomString(21),
		Password: randomPassword(),
	}
	err := suite.rr.users.Save(domain.User{ID: ap.UserID})
	suite.Require().NoError(err)
	err = suite.rr.authnPassword.Save(ap)
	suite.Require().NoError(err)

	return ap
}

func (suite *servicesTestSuite) Test_AuthnPassword_Login() {
	suite.Run("Login должен быть валидным", func() {
		out, err := suite.ss.authnPassword.Login(AuthnPasswordLoginInput{
			Login:    " inv ald login",
			Password: "somePassword123!",
		})
		suite.ErrorIs(err, ErrInvalidLogin)
		suite.Zero(out)
	})

	suite.Run("Password должен быть валидным", func() {
		out, err := suite.ss.authnPassword.Login(AuthnPasswordLoginInput{
			Login:    "someLogin",
			Password: "invalidpassword",
		})
		suite.ErrorIs(err, ErrInvalidPassword)
		suite.Zero(out)
	})

	suite.Run("неверные данные", func() {
		out, err := suite.ss.authnPassword.Login(AuthnPasswordLoginInput{
			Login:    "wrongLogin",
			Password: "wrongPassword123!",
		})
		suite.ErrorIs(err, ErrLoginOrPasswordDoesNotMatch)
		suite.Zero(out)
	})

	suite.Run("вернется Verified сессия", func() {
		// Создаем нового пользователя с AuthnPassword
		uwp := suite.newRndUserWithAuthnPassword()
		// Входим сессию с правильными данными
		input := AuthnPasswordLoginInput{
			Login:    uwp.Login,
			Password: uwp.Password,
		}
		output, err := suite.ss.authnPassword.Login(input)
		suite.NoError(err)
		suite.Require().NotZero(output)
		suite.Require().Equal(uwp.UserID, output.Session.UserID)
		suite.Require().Equal(domain.SessionStatusVerified, output.Session.Status)

		// Проверяем, что сессия сохранена в репозитории
		sessions, err := suite.rr.sessions.List(domain.SessionsFilter{
			Token: output.Session.Token,
		})
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.Equal(output.Session, sessions[0])
	})
}

func (suite *servicesTestSuite) Test_AuthnPassword_Registration() {
	suite.Run("Login должен быть валидным", func() {
		// Регистрация по логину паролю
		input := AuthnPasswordRegistrationInput{
			Login:    "",
			Password: randomPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.authnPassword.Registration(input)
		suite.ErrorIs(err, ErrInvalidLogin)
		suite.Zero(out)
	})

	suite.Run("Password должен быть валидным", func() {
		// Регистрация по логину паролю
		input := AuthnPasswordRegistrationInput{
			Login:    "login",
			Password: "",
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.authnPassword.Registration(input)
		suite.ErrorIs(err, ErrInvalidPassword)
		suite.Zero(out)
	})

	suite.Run("Name должен быть валидным", func() {
		// Регистрация по логину паролю
		input := AuthnPasswordRegistrationInput{
			Login:    "login",
			Password: randomPassword(),
			Name:     "",

			Nick: "nick",
		}
		out, err := suite.ss.authnPassword.Registration(input)
		suite.ErrorIs(err, ErrInvalidName)
		suite.Zero(out)
	})

	suite.Run("нельзя создать пользователя с существующим логином", func() {
		// Регистрация по логину паролю
		input := AuthnPasswordRegistrationInput{
			Login:    "login",
			Password: randomPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.authnPassword.Registration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Регистрация второй раз с существующим логином
		input = AuthnPasswordRegistrationInput{
			Login:    "login",
			Password: randomPassword(),
			Name:     "name2",
			Nick:     "nick2",
		}
		out, err = suite.ss.authnPassword.Registration(input)
		suite.ErrorIs(err, ErrLoginIsAlreadyInUse)
		suite.Zero(out)
	})

	suite.Run("после регистрации будет создан пользователь", func() {
		// Регистрация по логину паролю
		input := AuthnPasswordRegistrationInput{
			Login:    "login",
			Password: randomPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.authnPassword.Registration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		suite.Equal(input.Name, out.User.Name)
		suite.Equal(input.Nick, out.User.Nick)

		// Прочитать пользователя из репозитория
		users, err := suite.rr.users.List(domain.UsersFilter{})
		suite.Require().NoError(err)
		suite.Require().Len(users, 1)
		suite.Equal(out.User, users[0])
	})

	suite.Run("после регистрации будет создана сессия", func() {
		// Регистрация по логину паролю
		input := AuthnPasswordRegistrationInput{
			Login:    "login",
			Password: randomPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.authnPassword.Registration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Проверить сессию
		sessions, err := suite.rr.sessions.List(domain.SessionsFilter{})
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.Equal(out.Session, sessions[0])
		suite.Equal(domain.SessionStatusVerified, sessions[0].Status)
	})

	suite.Run("после регистрации будет создан метод входа", func() {
		// Регистрация по логину паролю
		input := AuthnPasswordRegistrationInput{
			Login:    "login",
			Password: randomPassword(),
			Name:     "name",
			Nick:     "nick",
		}
		out, err := suite.ss.authnPassword.Registration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Прочитать метод входа из репозитория
		aps, err := suite.rr.authnPassword.List(domain.AuthnPasswordFilter{Login: input.Login})
		suite.Require().NoError(err)
		suite.Require().Len(aps, 1)
		suite.Equal(out.User.ID, aps[0].UserID)
		suite.Equal(input.Login, aps[0].Login)
		suite.Equal(input.Password, aps[0].Password)
	})
}
