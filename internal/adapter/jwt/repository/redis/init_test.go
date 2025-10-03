package redisRegistry_test

import (
	redisRegistry "github.com/nice-pea/npchat/internal/adapter/jwt/repository/redis"
)

// Test_Init - тестовый сценарий для проверки функции инициализации Redis-клиента
func (suite *testSuite) Test_Init() {
	suite.Run("подключение к redis с валидной конфигурацией", func() {
		cfg := redisRegistry.Config{
			DSN: suite.DSN,
		}
		cli, err := redisRegistry.Init(cfg)

		suite.Require().NoError(err)
		suite.Assert().NotNil(cli)
		defer func() { _ = cli.Close() }()
	})
	suite.Run("подключение к redis с пустой конфигурацией", func() {
		cfg := redisRegistry.Config{}
		cli, err := redisRegistry.Init(cfg)

		suite.Require().Error(err)
		suite.Assert().Nil(cli)
	})
	suite.Run("подключение к redis с невалидной конфигурацией", func() {
		cfg := redisRegistry.Config{
			DSN: "241421.46334.14241.61253:253532325",
		}
		cli, err := redisRegistry.Init(cfg)

		suite.Require().Error(err)
		suite.Assert().Nil(cli)
	})
}
