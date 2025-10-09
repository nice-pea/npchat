package oauthComplete

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	testifySuite "github.com/stretchr/testify/suite"
	"golang.org/x/exp/maps"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	mockSessionn "github.com/nice-pea/npchat/internal/domain/sessionn/mocks"
	"github.com/nice-pea/npchat/internal/domain/userr"
	mockUserr "github.com/nice-pea/npchat/internal/domain/userr/mocks"
	serviceSuite "github.com/nice-pea/npchat/internal/usecases/suite"
	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
	mockOauth "github.com/nice-pea/npchat/internal/usecases/users/oauth/mocks"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) Test_OauthComplete() {
	const providerName = "mock"
	// usecase.Providers.Add(suite.Adapters.Oauth)

	suite.Run("UserCode обязательное поле", func() {
		usecase, _, _, _ := newUsecase(suite)
		out, err := usecase.OauthComplete(In{
			UserCode: "",
			Provider: providerName,
		})
		suite.ErrorIs(err, ErrInvalidUserCode)
		suite.Zero(out)
	})

	suite.Run("Provider обязательное поле", func() {
		usecase, _, _, _ := newUsecase(suite)
		out, err := usecase.OauthComplete(In{
			UserCode: uuid.NewString(),
			Provider: "",
		})
		suite.ErrorIs(err, ErrInvalidProvider)
		suite.Zero(out)
	})

	suite.Run("ошибка если у провайдера не совпадет UserCode", func() {
		usecase, _, _, _ := newUsecase(suite)
		usecase.Providers.Add(suite.Adapters.Oauth)
		out, err := usecase.OauthComplete(In{
			UserCode: uuid.NewString(),
			Provider: providerName,
		})
		suite.ErrorIs(err, ErrWrongUserCode)
		suite.Zero(out)
	})

	suite.Run("Provider должен быть известен в сервисе", func() {
		usecase, _, _, _ := newUsecase(suite)
		usecase.Providers.Add(suite.Adapters.Oauth)
		// Завершить регистрацию
		out, err := usecase.OauthComplete(In{
			UserCode: maps.Keys(suite.MockOauthTokens)[0],
			Provider: "unknownProvider",
		})
		suite.ErrorIs(err, oauth.ErrUnknownOauthProvider)
		suite.Zero(out)
	})

	suite.Run("если пользователя не существует, он будет создан", func() {
		usecase, _, mockRepoUsers, mockRepoSessions := newUsecase(suite)
		usecase.Providers.Add(suite.Adapters.Oauth)
		// Настройка моков
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()
		mockRepoUsers.EXPECT().InTransaction(mock.Anything).RunAndReturn(func(fn func(userr.Repository) error) error {
			return fn(mockRepoUsers)
		}).Once()
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()
		mockRepoUsers.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		mockRepoSessions.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		// Завершить регистрацию
		out, err := usecase.OauthComplete(In{
			UserCode: maps.Keys(suite.MockOauthTokens)[0],
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.NoError(err)
		suite.Require().NotZero(out)
	})

	suite.Run("новый пользователь будет иметь метод авторизации", func() {
		usecase, _, mockRepoUsers, mockRepoSessions := newUsecase(suite)
		usecase.Providers.Add(suite.Adapters.Oauth)

		pCode := maps.Keys(suite.MockOauthTokens)[0]
		pToken := suite.MockOauthTokens[pCode]
		pUser := suite.MockOauthUsers[pToken]

		input := In{
			UserCode: pCode,
			Provider: suite.Adapters.Oauth.Name(),
		}
		// Настройка моков
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()
		mockRepoUsers.EXPECT().InTransaction(mock.Anything).RunAndReturn(func(fn func(userr.Repository) error) error {
			return fn(mockRepoUsers)
		}).Once()
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()

		var repoUser userr.User
		mockRepoUsers.EXPECT().Upsert(mock.Anything).Run(func(user userr.User) {
			repoUser = user
		}).Return(nil).Once()
		mockRepoSessions.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		// Завершить регистрацию
		out, err := usecase.OauthComplete(input)
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить созданную связь
		suite.Require().Len(repoUser.OpenAuthUsers, 1)
		oauthUser := repoUser.OpenAuthUsers[0]
		suite.Equal(pUser.ID, oauthUser.ID)
		suite.Equal(input.Provider, oauthUser.Provider)
	})

	suite.Run("для нового пользователя будет создана сессия", func() {
		usecase, _, mockRepoUsers, mockRepoSessions := newUsecase(suite)
		usecase.Providers.Add(suite.Adapters.Oauth)
		// Настройка моков
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()
		mockRepoUsers.EXPECT().InTransaction(mock.Anything).RunAndReturn(func(fn func(userr.Repository) error) error {
			return fn(mockRepoUsers)
		}).Once()
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()
		mockRepoUsers.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		var repoSession sessionn.Session
		mockRepoSessions.EXPECT().Upsert(mock.Anything).Run(func(session sessionn.Session) {
			repoSession = session
		}).Return(nil).Once()
		// Завершить регистрацию
		out, err := usecase.OauthComplete(In{
			UserCode: maps.Keys(suite.MockOauthTokens)[0],
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить сессию

		suite.EqualSessions(out.Session, repoSession)
		suite.Equal(sessionn.StatusVerified, repoSession.Status)
	})

	suite.Run("последующие попытки авторизации того же пользователя, будут создавать новые сессии", func() {
		usecase, _, mockRepoUsers, mockRepoSessions := newUsecase(suite)
		usecase.Providers.Add(suite.Adapters.Oauth)

		pCode := maps.Keys(suite.MockOauthTokens)[0]

		var sessions []sessionn.Session
		var savedUser userr.User
		// Настройка моков
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()
		mockRepoUsers.EXPECT().InTransaction(mock.Anything).RunAndReturn(func(fn func(userr.Repository) error) error {
			return fn(mockRepoUsers)
		}).Once()
		mockRepoUsers.EXPECT().List(mock.Anything).Return(nil, nil).Once()
		mockRepoUsers.EXPECT().Upsert(mock.Anything).Run(func(user userr.User) {
			savedUser = user
		}).Return(nil).Once()
		mockRepoSessions.EXPECT().Upsert(mock.Anything).Run(func(session sessionn.Session) {
			sessions = append(sessions, session)
		}).Return(nil).Once()
		// Завершить регистрацию
		out1, err := usecase.OauthComplete(In{
			UserCode: pCode,
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.Require().NoError(err)
		suite.NotZero(out1)
		// Настройка моков
		mockRepoUsers.EXPECT().InTransaction(mock.Anything).RunAndReturn(func(fn func(userr.Repository) error) error {
			return fn(mockRepoUsers)
		}).Once()
		mockRepoUsers.EXPECT().List(mock.Anything).Return([]userr.User{savedUser}, nil).Once()
		mockRepoUsers.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		mockRepoSessions.EXPECT().Upsert(mock.Anything).Run(func(session sessionn.Session) {
			sessions = append(sessions, session)
		}).Return(nil).Once()
		// Завершить регистрацию, с UserCode связанным с тем же пользователем провайдера
		out2, err := usecase.OauthComplete(In{
			UserCode: pCode,
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.Require().NoError(err)
		suite.NotZero(out2)

		suite.Require().Len(sessions, 2)
		// После первого входа (регистрации)
		suite.EqualSessions(out1.Session, sessions[0])
		suite.Equal(sessionn.StatusVerified, sessions[0].Status)
		// После второго входа
		suite.EqualSessions(out2.Session, sessions[1])
		suite.Equal(sessionn.StatusVerified, sessions[1].Status)
	})
}

func newUsecase(suite *testSuite) (*OauthCompleteUsecase, *mockOauth.Provider, *mockUserr.Repository, *mockSessionn.Repository) {
	usecase := &OauthCompleteUsecase{
		Providers:    oauth.Providers{},
		Repo:         suite.RR.Users,
		SessionsRepo: suite.RR.Sessions,
	}
	mockUsers := suite.RR.Users
	mockSessions := suite.RR.Sessions
	mockOauth := suite.Adapters.Oauth
	return usecase, mockOauth, mockUsers, mockSessions
}
