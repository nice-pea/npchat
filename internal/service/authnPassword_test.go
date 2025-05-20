package service

import (
	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func (suite *servicesTestSuite) newRndUserWithAuthnPassword() domain.AuthnPassword {
	ap := domain.AuthnPassword{
		UserID:   uuid.NewString(),
		Login:    uuid.NewString(),
		Password: uuid.NewString(),
	}
	err := suite.rr.users.Save(domain.User{ID: ap.UserID})
	suite.Require().NoError(err)
	err = suite.rr.authnPassword.Save(ap)
	suite.Require().NoError(err)

	return ap
}

func (suite *servicesTestSuite) Test_AuthnPassword_Login() {
	suite.Run("неверные данные", func() {
		session, err := suite.ss.authnPassword.Login(AuthnPasswordLoginInput{
			Login:    "wronglogin",
			Password: "wrongpassword",
		})
		suite.Error(err)
		suite.Zero(session)
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
	suite.Run("Login обязательное поле", func() {
		// Регистрация по логину паролю
		input := AuthnPasswordRegistrationInput{
			Login:    "",
			Password: "Password123!",
			Name:     "name",
			Nick:     "nick",
		}
		user, err := suite.ss.authnPassword.Registration(input)
		suite.ErrorIs(err, ErrInvalidLogin)
		suite.Zero(user)
	})

	suite.Run("Password обязательное поле", func() {
		// Регистрация по логину паролю
		input := AuthnPasswordRegistrationInput{
			Login:    "login",
			Password: "",
			Name:     "name",
			Nick:     "nick",
		}
		user, err := suite.ss.authnPassword.Registration(input)
		suite.ErrorIs(err, ErrInvalidPassword)
		suite.Zero(user)
	})

	suite.Run("Name обязательное поле", func() {
		// Регистрация по логину паролю
		input := AuthnPasswordRegistrationInput{
			Login:    "login",
			Password: "Password123!",
			Name:     "",
			Nick:     "nick",
		}
		user, err := suite.ss.authnPassword.Registration(input)
		suite.ErrorIs(err, ErrInvalidName)
		suite.Zero(user)
	})

	suite.Run("пользователя и метод можно прочитать из репозитория", func() {
		// Регистрация по логину паролю
		input := AuthnPasswordRegistrationInput{
			Login:    "login",
			Password: "Password123!",
			Name:     "name",
			Nick:     "nick",
		}
		user, err := suite.ss.authnPassword.Registration(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(user)
		suite.Equal(input.Name, user.Name)
		suite.Equal(input.Nick, user.Nick)

		// Прочитать пользователя из репозитория
		users, err := suite.rr.users.List(domain.UsersFilter{})
		suite.Require().NoError(err)
		suite.Require().Len(users, 1)
		suite.Equal(user, users[0])

		// Прочитать метод входа из репозитория
		aps, err := suite.rr.authnPassword.List(domain.AuthnPasswordFilter{Login: input.Login})
		suite.Require().NoError(err)
		suite.Require().Len(aps, 1)
		suite.Equal(user.ID, aps[0].UserID)
		suite.Equal(input.Login, aps[0].Login)
		suite.Equal(input.Password, aps[0].Password)
	})
}
