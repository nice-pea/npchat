package initOAuthRegistration

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

func (suite *testSuite) Test_OAuth_InitRegistration() {
	usecase := &InitOAuthRegistrationUsecase{
		Providers: oauth.OAuthProviders{},
	}
	usecase.Providers.Add(suite.Adapters.Oauth)

	suite.Run("Provider обязательное поле", func() {
		// Инициализация регистрации
		out, err := usecase.InitOAuthRegistration(In{
			Provider: "",
		})
		suite.ErrorIs(err, ErrInvalidProvider)
		suite.Zero(out)
	})

	suite.Run("Provider должен быть известен в сервисе", func() {
		// Инициализация регистрации
		input := In{
			Provider: "unknownProvider",
		}
		out, err := usecase.InitOAuthRegistration(input)
		suite.ErrorIs(err, oauth.ErrUnknownOAuthProvider)
		suite.Zero(out)
	})

	suite.Run("инициализация вернет валидный url", func() {
		// Инициализация регистрации
		out, err := usecase.InitOAuthRegistration(In{
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
}
