package completeOauthRegistration

import (
	"testing"

	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"
	"golang.org/x/exp/maps"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	serviceSuite "github.com/nice-pea/npchat/internal/usecases/suite"
	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) Test_Oauth_CompleteRegistration() {
	usecase := &CompleteOauthRegistrationUsecase{
		Repo:         suite.RR.Users,
		SessionsRepo: suite.RR.Sessions,
		Providers:    oauth.OauthProviders{},
	}
	usecase.Providers.Add(suite.Adapters.Oauth)

	suite.Run("UserCode обязательное поле", func() {
		input := In{
			UserCode: "",
			Provider: suite.Adapters.Oauth.Name(),
		}
		out, err := usecase.CompleteOauthRegistration(input)
		suite.ErrorIs(err, ErrInvalidUserCode)
		suite.Zero(out)
	})

	suite.Run("Provider обязательное поле", func() {
		input := In{
			UserCode: uuid.NewString(),
			Provider: "",
		}
		out, err := usecase.CompleteOauthRegistration(input)
		suite.ErrorIs(err, ErrInvalidProvider)
		suite.Zero(out)
	})

	suite.Run("ошибка если у провайдера не совпадет UserCode", func() {
		input := In{
			UserCode: uuid.NewString(),
			Provider: suite.Adapters.Oauth.Name(),
		}
		out, err := usecase.CompleteOauthRegistration(input)
		suite.ErrorIs(err, ErrWrongUserCode)
		suite.Zero(out)
	})

	suite.Run("Provider должен быть известен в сервисе", func() {
		// Завершить регистрацию
		input := In{
			UserCode: maps.Keys(suite.MockOauthTokens)[0],
			Provider: "unknownProvider",
		}
		out, err := usecase.CompleteOauthRegistration(input)
		suite.ErrorIs(err, oauth.ErrUnknownOauthProvider)
		suite.Zero(out)
	})

	suite.Run("после регистрации будет создан пользователь", func() {
		// Завершить регистрацию
		input := In{
			UserCode: maps.Keys(suite.MockOauthTokens)[0],
			Provider: suite.Adapters.Oauth.Name(),
		}
		out, err := usecase.CompleteOauthRegistration(input)
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить users репозиторий
		users, err := suite.RR.Users.List(userr.Filter{})
		suite.NoError(err)
		suite.Require().Len(users, 1)
		suite.True(out.User.Equal(users[0]))
	})

	suite.Run("после регистрации будет создан метод авторизации", func() {
		pCode := maps.Keys(suite.MockOauthTokens)[0]
		pToken := suite.MockOauthTokens[pCode]
		pUser := suite.MockOauthUsers[pToken]

		// Завершить регистрацию
		input := In{
			UserCode: pCode,
			Provider: suite.Adapters.Oauth.Name(),
		}
		out, err := usecase.CompleteOauthRegistration(input)
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить созданную связь
		users, err := suite.RR.Users.List(userr.Filter{})
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
		input := In{
			UserCode: maps.Keys(suite.MockOauthTokens)[0],
			Provider: suite.Adapters.Oauth.Name(),
		}
		out, err := usecase.CompleteOauthRegistration(input)
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить сессию
		sessions, err := suite.RR.Sessions.List(sessionn.Filter{})
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.EqualSessions(out.Session, sessions[0])
		suite.Equal(sessionn.StatusVerified, sessions[0].Status)
	})

	suite.Run("невозможно дважды зарегистрироваться на одного пользователя провайдера", func() {
		pCode := maps.Keys(suite.MockOauthTokens)[0]

		// Завершить регистрацию
		input := In{
			UserCode: pCode,
			Provider: suite.Adapters.Oauth.Name(),
		}
		out, err := usecase.CompleteOauthRegistration(input)
		suite.Require().NoError(err)
		suite.NotZero(out)

		// Завершить регистрацию, с UserCode связанным с тем же пользователем провайдера
		input = In{
			UserCode: pCode,
			Provider: suite.Adapters.Oauth.Name(),
		}
		out, err = usecase.CompleteOauthRegistration(input)
		suite.Error(err, ErrProvidersUserIsAlreadyLinked)
		suite.Zero(out)
	})
}
