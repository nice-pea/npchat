package service

import (
	"net/url"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"

	"github.com/saime-0/nice-pea-chat/internal/domain/sessionn"
	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

func (suite *servicesTestSuite) Test_OAuth_InitRegistration() {
	suite.Run("Provider обязательное поле", func() {
		// Инициализация регистрации
		out, err := suite.ss.users.InitOAuthRegistration(InitOAuthRegistrationIn{
			Provider: "",
		})
		suite.ErrorIs(err, ErrInvalidProvider)
		suite.Zero(out)
	})

	suite.Run("Provider должен быть известен в сервисе", func() {
		// Инициализация регистрации
		input := InitOAuthRegistrationIn{
			Provider: "unknownProvider",
		}
		out, err := suite.ss.users.InitOAuthRegistration(input)
		suite.ErrorIs(err, ErrUnknownOAuthProvider)
		suite.Zero(out)
	})

	suite.Run("инициализация вернет валидный url", func() {
		// Инициализация регистрации
		out, err := suite.ss.users.InitOAuthRegistration(InitOAuthRegistrationIn{
			Provider: suite.ad.oauth.Name(),
		})
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Вернется url
		parsedUrl, err := url.Parse(out.RedirectURL)
		suite.NoError(err)

		// Есть query-параметр state
		state := parsedUrl.Query().Get("state")
		suite.NotZero(state)
		// Есть query-параметр state
		code := parsedUrl.Query().Get("code")
		suite.NotZero(code)
	})
}

func (suite *servicesTestSuite) Test_OAuth_CompleteRegistration() {
	suite.Run("UserCode обязательное поле", func() {
		input := CompeteOAuthRegistrationIn{
			UserCode: "",
			Provider: suite.ad.oauth.Name(),
		}
		out, err := suite.ss.users.CompeteOAuthRegistration(input)
		suite.ErrorIs(err, ErrInvalidUserCode)
		suite.Zero(out)
	})

	suite.Run("Provider обязательное поле", func() {
		input := CompeteOAuthRegistrationIn{
			UserCode: uuid.NewString(),
			Provider: "",
		}
		out, err := suite.ss.users.CompeteOAuthRegistration(input)
		suite.ErrorIs(err, ErrInvalidProvider)
		suite.Zero(out)
	})

	suite.Run("ошибка если у провайдера не совпадет UserCode", func() {
		input := CompeteOAuthRegistrationIn{
			UserCode: uuid.NewString(),
			Provider: suite.ad.oauth.Name(),
		}
		out, err := suite.ss.users.CompeteOAuthRegistration(input)
		suite.ErrorIs(err, ErrWrongUserCode)
		suite.Zero(out)
	})

	suite.Run("Provider должен быть известен в сервисе", func() {
		// Завершить регистрацию
		input := CompeteOAuthRegistrationIn{
			UserCode: maps.Keys(suite.mockOAuthTokens)[0],
			Provider: "unknownProvider",
		}
		out, err := suite.ss.users.CompeteOAuthRegistration(input)
		suite.ErrorIs(err, ErrUnknownOAuthProvider)
		suite.Zero(out)
	})

	suite.Run("после регистрации будет создан пользователь", func() {
		// Завершить регистрацию
		input := CompeteOAuthRegistrationIn{
			UserCode: maps.Keys(suite.mockOAuthTokens)[0],
			Provider: suite.ad.oauth.Name(),
		}
		out, err := suite.ss.users.CompeteOAuthRegistration(input)
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить users репозиторий
		users, err := suite.rr.users.List(userr.Filter{})
		suite.NoError(err)
		suite.Require().Len(users, 1)
		suite.Equal(out.User, users[0])
	})

	suite.Run("после регистрации будет создан метод авторизации", func() {
		pCode := maps.Keys(suite.mockOAuthTokens)[0]
		pToken := suite.mockOAuthTokens[pCode]
		pUser := suite.mockOAuthUsers[pToken]

		// Завершить регистрацию
		input := CompeteOAuthRegistrationIn{
			UserCode: pCode,
			Provider: suite.ad.oauth.Name(),
		}
		out, err := suite.ss.users.CompeteOAuthRegistration(input)
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить созданную связь
		users, err := suite.rr.users.List(userr.Filter{})
		suite.NoError(err)
		suite.Require().Len(users, 1)
		user := users[0]
		suite.Require().Len(user.OpenAuthUsers, 1)
		oauthUser := user.OpenAuthUsers[0]
		suite.Equal(pUser.ID, oauthUser.ID)
		suite.Equal(input.Provider, oauthUser.Provider)
	})

	suite.Run("после регистрации будет создана сессия", func() {
		// Завершить регистрацию
		input := CompeteOAuthRegistrationIn{
			UserCode: maps.Keys(suite.mockOAuthTokens)[0],
			Provider: suite.ad.oauth.Name(),
		}
		out, err := suite.ss.users.CompeteOAuthRegistration(input)
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить сессию
		sessions, err := suite.rr.sessions.List(sessionn.Filter{})
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.Equal(out.Session, sessions[0])
		suite.Equal(sessionn.StatusVerified, sessions[0].Status)
	})

	suite.Run("невозможно дважды зарегистрироваться на одного пользователя провайдера", func() {
		pCode := maps.Keys(suite.mockOAuthTokens)[0]

		// Завершить регистрацию
		input := CompeteOAuthRegistrationIn{
			UserCode: pCode,
			Provider: suite.ad.oauth.Name(),
		}
		out, err := suite.ss.users.CompeteOAuthRegistration(input)
		suite.Require().NoError(err)
		suite.NotZero(out)

		// Завершить регистрацию, с UserCode связанным с тем же пользователем провайдера
		input = CompeteOAuthRegistrationIn{
			UserCode: pCode,
			Provider: suite.ad.oauth.Name(),
		}
		out, err = suite.ss.users.CompeteOAuthRegistration(input)
		suite.Error(err, ErrProvidersUserIsAlreadyLinked)
		suite.Zero(out)
	})
}
