package service

import (
	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func (suite *servicesTestSuite) newRndUserWithLoginCredentials() domain.LoginCredentials {
	credentials := domain.LoginCredentials{
		UserID:   uuid.NewString(),
		Login:    uuid.NewString(),
		Password: uuid.NewString(),
	}
	err := suite.rr.users.Save(domain.User{ID: credentials.UserID})
	suite.Require().NoError(err)
	err = suite.rr.loginCredentials.Save(credentials)
	suite.Require().NoError(err)

	return credentials
}

func (suite *servicesTestSuite) Test_LoginCredentials_Login() {
	suite.Run("неверные данные", func() {
		session, err := suite.ss.loginCredentials.Login(LoginByCredentialsInput{
			Login:    "wronglogin",
			Password: "wrongpassword",
		})
		suite.Error(err)
		suite.Zero(session)
	})
	suite.Run("вернется Verified сессия", func() {
		// Создаем нового пользователя с login credentials
		uwlc := suite.newRndUserWithLoginCredentials()
		// Входим сессию с правильными данными
		input := LoginByCredentialsInput{
			Login:    uwlc.Login,
			Password: uwlc.Password,
		}
		session, err := suite.ss.loginCredentials.Login(input)
		suite.NoError(err)
		suite.Require().NotZero(session)
		suite.Require().Equal(uwlc.UserID, session.UserID)
		suite.Require().Equal(domain.SessionStatusVerified, session.Status)

		// Проверяем, что сессия сохранена в репозитории
		sessions, err := suite.rr.sessions.List(domain.SessionsFilter{
			Token: session.Token,
		})
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.Equal(session, sessions[0])
	})
}
