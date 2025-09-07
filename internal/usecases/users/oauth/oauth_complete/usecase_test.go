package oauthComplete

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

func (suite *testSuite) Test_OauthComplete() {
	usecase := &OauthCompleteUsecase{
		Repo:         suite.RR.Users,
		SessionsRepo: suite.RR.Sessions,
		Providers:    oauth.Providers{},
	}
	usecase.Providers.Add(suite.Adapters.Oauth)

	suite.Run("UserCode обязательное поле", func() {
		out, err := usecase.OauthComplete(In{
			UserCode: "",
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.ErrorIs(err, ErrInvalidUserCode)
		suite.Zero(out)
	})

	suite.Run("Provider обязательное поле", func() {
		out, err := usecase.OauthComplete(In{
			UserCode: uuid.NewString(),
			Provider: "",
		})
		suite.ErrorIs(err, ErrInvalidProvider)
		suite.Zero(out)
	})

	suite.Run("ошибка если у провайдера не совпадет UserCode", func() {
		out, err := usecase.OauthComplete(In{
			UserCode: uuid.NewString(),
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.ErrorIs(err, ErrWrongUserCode)
		suite.Zero(out)
	})

	suite.Run("Provider должен быть известен в сервисе", func() {
		// Завершить регистрацию
		out, err := usecase.OauthComplete(In{
			UserCode: maps.Keys(suite.MockOauthTokens)[0],
			Provider: "unknownProvider",
		})
		suite.ErrorIs(err, oauth.ErrUnknownOauthProvider)
		suite.Zero(out)
	})

	suite.Run("если пользователя не существует, он будет создан", func() {
		// Завершить регистрацию
		out, err := usecase.OauthComplete(In{
			UserCode: maps.Keys(suite.MockOauthTokens)[0],
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить users репозиторий
		users, err := suite.RR.Users.List(userr.Filter{})
		suite.NoError(err)
		suite.Require().Len(users, 1)
		suite.True(out.User.Equal(users[0]))
	})

	suite.Run("новый пользователь будет иметь метод авторизации", func() {
		pCode := maps.Keys(suite.MockOauthTokens)[0]
		pToken := suite.MockOauthTokens[pCode]
		pUser := suite.MockOauthUsers[pToken]

		// Завершить регистрацию
		input := In{
			UserCode: pCode,
			Provider: suite.Adapters.Oauth.Name(),
		}
		out, err := usecase.OauthComplete(input)
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

	suite.Run("для нового пользователя будет создана сессия", func() {
		// Завершить регистрацию
		out, err := usecase.OauthComplete(In{
			UserCode: maps.Keys(suite.MockOauthTokens)[0],
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить сессию
		sessions, err := suite.RR.Sessions.List(sessionn.Filter{})
		suite.NoError(err)
		suite.Require().Len(sessions, 1)
		suite.EqualSessions(out.Session, sessions[0])
		suite.Equal(sessionn.StatusVerified, sessions[0].Status)
	})

	suite.Run("последующие попытки авторизации того же пользователя, будут создавать новые сессии", func() {
		pCode := maps.Keys(suite.MockOauthTokens)[0]

		// Завершить регистрацию
		out1, err := usecase.OauthComplete(In{
			UserCode: pCode,
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.Require().NoError(err)
		suite.NotZero(out1)

		// Завершить регистрацию, с UserCode связанным с тем же пользователем провайдера
		out2, err := usecase.OauthComplete(In{
			UserCode: pCode,
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.Require().NoError(err)
		suite.NotZero(out2)

		sessions, err := suite.RR.Sessions.List(sessionn.Filter{})
		suite.NoError(err)
		suite.Require().Len(sessions, 2)
		// После первого входа (регистрации)
		suite.EqualSessions(out1.Session, sessions[0])
		suite.Equal(sessionn.StatusVerified, sessions[0].Status)
		// После второго входа
		suite.EqualSessions(out2.Session, sessions[1])
		suite.Equal(sessionn.StatusVerified, sessions[1].Status)
	})
}
