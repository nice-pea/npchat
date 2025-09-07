package oauthAuthorize

import (
	"net/url"
	"testing"

	testifySuite "github.com/stretchr/testify/suite"

	serviceSuite "github.com/nice-pea/npchat/internal/usecases/suite"
	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) Test_OauthAuthorize() {
	usecase := &OauthAuthorizeUsecase{
		Providers: oauth.Providers{},
	}
	usecase.Providers.Add(suite.Adapters.Oauth)

	suite.Run("Provider обязательные поля, должен быть известен в сервисе", func() {
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
		// Инициализация регистрации
		out, err := usecase.OauthAuthorize(In{
			Provider: suite.Adapters.Oauth.Name(),
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
		// Инициализация регистрации
		out, err := usecase.OauthAuthorize(In{
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.NoError(err)
		suite.Require().NotZero(out)
		suite.NotEmpty(out.State)

		// Повторная инициализация даёт новый state
		out2, err := usecase.OauthAuthorize(In{
			Provider: suite.Adapters.Oauth.Name(),
		})
		suite.NoError(err)
		suite.Require().NotZero(out2)
		suite.NotEmpty(out2.State)

		suite.NotEqual(out.State, out2.State)
	})
}
