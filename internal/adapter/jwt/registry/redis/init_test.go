package redisRegistry_test

import (
	redisRegistry "github.com/nice-pea/npchat/internal/adapter/jwt/registry/redis"
)

// Test_Init тестовый сценарий для проверки функции инициализации Redis-клиента
func (suite *testSuite) Test_Init() {
	suite.Run("подключение к redis с валидной конфигурацией", func() {
		cfg := redisRegistry.Config{
			DSN: suite.DSN,
		}
		cli, err := redisRegistry.Init(cfg)

		suite.Require().NoError(err)
		suite.NotNil(cli)
		err = cli.Close()
		suite.NoError(err)
	})

	suite.Run("подключение к redis с пустой конфигурацией", func() {
		cfg := redisRegistry.Config{}
		cli, err := redisRegistry.Init(cfg)

		suite.Require().Error(err)
		suite.Nil(cli)
	})

	suite.Run("подключение к redis с невалидной конфигурацией", func() {
		cfg := redisRegistry.Config{
			DSN: "241421.46334.14241.61253:253532325",
		}
		cli, err := redisRegistry.Init(cfg)

		suite.Require().Error(err)
		suite.Nil(cli)
	})
}
