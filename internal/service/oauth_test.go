package service

import (
	"net/url"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func (suite *servicesTestSuite) oauthRegistrationInit() OAuthRegistrationInitOut {
	out, err := suite.ss.oauth.InitRegistration(OAuthInitRegistrationInput{
		Provider: testProvider,
	})
	suite.Require().NoError(err)
	suite.Require().NotZero(out)

	_, err = url.Parse(out.RedirectURL)
	suite.Require().NoError(err)

	return out
}

func (suite *servicesTestSuite) oauthFullRegistration() domain.User {
	// Инициализация регистрации
	suite.oauthRegistrationInit()
	pCode := maps.Keys(suite.mockOAuthTokens)[0]

	// Завершить регистрацию
	input := OAuthCompeteRegistrationInput{
		UserCode: pCode,
		Provider: testProvider,
	}
	user, err := suite.ss.oauth.CompeteRegistration(input)
	suite.Require().NoError(err)

	return user
}

func (suite *servicesTestSuite) Test_OAuth_InitRegistration() {
	suite.Run("Provider обязательное поле", func() {
		// Инициализация регистрации
		out, err := suite.ss.oauth.InitRegistration(OAuthInitRegistrationInput{
			Provider: "",
		})
		suite.ErrorIs(err, ErrInvalidProvider)
		suite.Zero(out)
	})

	suite.Run("инициализация вернет валидный url", func() {
		// Инициализация регистрации
		out, err := suite.ss.oauth.InitRegistration(OAuthInitRegistrationInput{
			Provider: testProvider,
		})
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Вернется url
		parsedUrl, err := url.Parse(out.RedirectURL)
		suite.NoError(err)

		// Есть query-параметр state
		state := parsedUrl.Query().Get("state")
		suite.NotZero(state)
	})
}

func (suite *servicesTestSuite) Test_OAuth_CompleteRegistration() {
	suite.Run("UserCode обязательное поле", func() {
		input := OAuthCompeteRegistrationInput{
			UserCode: "",
			Provider: testProvider,
		}
		out, err := suite.ss.oauth.CompeteRegistration(input)
		suite.ErrorIs(err, ErrInvalidUserCode)
		suite.Zero(out)
	})

	suite.Run("Provider обязательное поле", func() {
		input := OAuthCompeteRegistrationInput{
			UserCode: uuid.NewString(),
			Provider: "",
		}
		out, err := suite.ss.oauth.CompeteRegistration(input)
		suite.ErrorIs(err, ErrInvalidProvider)
		suite.Zero(out)
	})

	suite.Run("неверный UserCode", func() {
		input := OAuthCompeteRegistrationInput{
			UserCode: uuid.NewString(),
			Provider: testProvider,
		}
		out, err := suite.ss.oauth.CompeteRegistration(input)
		suite.ErrorIs(err, ErrWrongUserCode)
		suite.Zero(out)
	})

	suite.Run("Provider должен быть известен в сервисе", func() {
		// Завершить регистрацию
		input := OAuthCompeteRegistrationInput{
			UserCode: maps.Keys(suite.mockOAuthTokens)[0],
			Provider: "unknownProvider",
		}
		out, err := suite.ss.oauth.CompeteRegistration(input)
		suite.ErrorIs(err, ErrUnknownOAuthProvider)
		suite.Zero(out)
	})

	suite.Run("после регистрации пользователя можно прочитать", func() {
		pCode := maps.Keys(suite.mockOAuthTokens)[0]
		pToken := suite.mockOAuthTokens[pCode]
		pUser := suite.mockOAuthUsers[pToken]

		// Завершить регистрацию
		input := OAuthCompeteRegistrationInput{
			UserCode: pCode,
			Provider: testProvider,
		}
		user, err := suite.ss.oauth.CompeteRegistration(input)
		suite.NoError(err)
		suite.Require().NotZero(user)

		// Проверить users репозиторий
		users, err := suite.rr.users.List(domain.UsersFilter{})
		suite.NoError(err)
		suite.Require().Len(users, 1)
		suite.Equal(user, users[0])
		// Поле провайдера
		suite.Equal(pUser.Name, user.Name)

		// Проверить oauth репозиторий
		links, err := suite.rr.oauth.ListLinks(domain.OAuthListLinksFilter{})
		suite.NoError(err)
		suite.Require().Len(links, 1)
		suite.Equal(pUser.ID, links[0].ExternalID)
		suite.Equal(user.ID, links[0].UserID)
		suite.Equal(testProvider, links[0].Provider)
	})

	suite.Run("невозможно дважды зарегистрироваться на одного пользователя провайдера", func() {
		pCode := maps.Keys(suite.mockOAuthTokens)[0]

		// Завершить регистрацию
		input := OAuthCompeteRegistrationInput{
			UserCode: pCode,
			Provider: testProvider,
		}
		user, err := suite.ss.oauth.CompeteRegistration(input)
		suite.Require().NoError(err)
		suite.NotZero(user)

		// Завершить регистрацию, с UserCode связанным с тем же пользователем провайдера
		input = OAuthCompeteRegistrationInput{
			UserCode: pCode,
			Provider: testProvider,
		}
		user, err = suite.ss.oauth.CompeteRegistration(input)
		suite.Error(err, ErrProvidersUserIsAlreadyLinked)
		suite.Zero(user)
	})
}
