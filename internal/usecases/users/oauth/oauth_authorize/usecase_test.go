package oauthAuthorize

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/mock"
	testifySuite "github.com/stretchr/testify/suite"

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

func (suite *testSuite) Test_OauthAuthorize() {
	suite.Run("Provider обязательные поля, должен быть известен в сервисе", func() {
		usecase, _ := newUsecase(suite)
		out, err := usecase.OauthAuthorize(In{
			Provider: "",
		})
		suite.ErrorIs(err, ErrInvalidProvider)
		suite.Zero(out)

		out, err = usecase.OauthAuthorize(In{
			Provider: "unknownProvider",
		})
		suite.ErrorIs(err, oauth.ErrUnknownOauthProvider)
		suite.Zero(out)
	})

	suite.Run("инициализация вернет валидный url", func() {
		usecase, mockOauth := newUsecase(suite)
		mockOauth.EXPECT().Name().Return("google.com").Once()
		usecase.Providers.Add(suite.Adapters.Oauth)
		// Инициализация регистрации
		mockOauth.EXPECT().Name().Return("google.com").Once()
		mockOauth.EXPECT().AuthorizationURL(mock.Anything).
			Return("https://accounts.google.com/o/oauth2/auth?state=STATE123&code=CODE456").
			Once()
		out, err := usecase.OauthAuthorize(In{
			Provider: mockOauth.Name(),
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

	suite.Run("инициализация вернет случайную строку в state", func() {
		usecase, mockOauth := newUsecase(suite)
		mockOauth.EXPECT().Name().Return("google.com").Once()
		usecase.Providers.Add(suite.Adapters.Oauth)
		// Инициализация регистрации
		mockOauth.EXPECT().Name().Return("google.com").Once()
		mockOauth.EXPECT().AuthorizationURL(mock.Anything).
			Return("https://accounts.google.com/o/oauth2/auth?state=STATE123&code=CODE456").
			Once()
		out, err := usecase.OauthAuthorize(In{
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.NoError(err)
		suite.Require().NotZero(out)
		suite.NotEmpty(out.State)

		// Повторная инициализация даёт новый state
		mockOauth.EXPECT().Name().Return("google.com").Once()
		mockOauth.EXPECT().AuthorizationURL(mock.Anything).
			Return("https://accounts.google.com/o/oauth2/auth?state=STATE234&code=CODE567").
			Once()
		out2, err := usecase.OauthAuthorize(In{
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.NoError(err)
		suite.Require().NotZero(out2)
		suite.NotEmpty(out2.State)

		suite.NotEqual(out.State, out2.State)
	})
}

func newUsecase(suite *testSuite) (*OauthAuthorizeUsecase, *mockOauth.Provider) {
	usecase := &OauthAuthorizeUsecase{
		Providers: oauth.Providers{},
	}
	mockOauth := suite.Adapters.Oauth
	return usecase, mockOauth
}
