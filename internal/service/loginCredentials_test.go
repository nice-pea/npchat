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
		session, err := suite.ss.authnPassword.Login(input)
		suite.NoError(err)
		suite.Require().NotZero(session)
		suite.Require().Equal(uwp.UserID, session.UserID)
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
