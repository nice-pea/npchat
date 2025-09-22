package redisCache_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type testSuite struct {
	suite.Suite
	cleanUp     func()
	ExposedAddr string
}

func Test_TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) newRedisContainer() {

	req := testcontainers.ContainerRequest{
		Image:        "redis:8.2.1",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	suite.Require().NoError(err)

	suite.cleanUp = func() {
		suite.Require().NotNil(container)
		container.Terminate(ctx)
	}

	// Получаем адрес контейнера
	endpoint, err := container.Endpoint(ctx, "")
	suite.Require().NoError(err)

	suite.ExposedAddr = endpoint

}

// SetupSuite выполняется один раз перед всеми тестами
func (suite *testSuite) SetupSuite() {
	suite.newRedisContainer()
}

// SetupTest выполняется перед каждым тестом
func (suite *testSuite) SetupTest() {
	println("SetupTest executed")

}

// TearDownTest выполняется после каждого теста
func (suite *testSuite) TearDownTest() {
	println("TearDownTest executed")
}

// TearDownSuite выполняется после всех тестов
func (suite *testSuite) TearDownSuite() {
	suite.cleanUp()
}
