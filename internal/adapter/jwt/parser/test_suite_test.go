package jwtParser

import (
	"testing"

	jwt2 "github.com/cristalhq/jwt/v5"
	"github.com/nice-pea/npchat/internal/adapter/jwt"
	"github.com/stretchr/testify/suite"

	mockJwtParser "github.com/nice-pea/npchat/internal/adapter/jwt/parser/mocks"
)

type testSuite struct {
	suite.Suite
	cfg          jwt.Config
	registryMock *mockJwtParser.Registry
	Parser       Parser
}

// Test_TestSuite запускает набор тестов
func Test_TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

// SetupSuite выполняется один раз перед всеми тестами
func (suite *testSuite) SetupSuite() {

	suite.cfg = jwt.Config{
		SecretKey:                   "secret",
		VerifyTokenWithInvalidation: true,
	}

}

// SetupSubTest выполняется перед каждым подтестом
func (suite *testSuite) SetupSubTest() {
	// создаем mockRegistry
	suite.registryMock = mockJwtParser.NewRegistry(suite.T())

	// создаем Parser
	suite.cfg.VerifyTokenWithInvalidation = true

	suite.Parser = Parser{
		Config:   suite.cfg,
		Registry: suite.registryMock,
	}
}

// createJWT создает JWT токен для тестов
func (suite *testSuite) createJWT(secret string, claims map[string]any) string {
	// создаем Signer
	signer, err := jwt2.NewSignerHS(jwt2.HS256, []byte(secret))
	suite.Require().NoError(err)
	// создаем Builder
	builder := jwt2.NewBuilder(signer)

	// создаем токен
	token, err := builder.Build(claims)
	suite.Require().NoError(err)

	return token.String()
}

// parserWithOutRegistry создает парсер без Registry
func (suite *testSuite) parserWithOutRegistry(secret string) Parser {
	return Parser{Config: jwt.Config{SecretKey: secret}}
}
