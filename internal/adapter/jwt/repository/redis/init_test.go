package redisCache_test

import (
	redisCache "github.com/nice-pea/npchat/internal/adapter/jwt/repository/redis"
)

func (suite *testSuite) Test_Init() {
	suite.Run("подключение к redis с неверной конфигурацией", func() {
		cfg := redisCache.Config{
			Addr: suite.ExposedAddr,
		}
		cli, err := redisCache.Init(cfg)

		suite.Require().NoError(err)
		suite.Assert().NotNil(cli)

	})
	suite.Run("подключение к redis с пустой конфигурацией", func() {
		cfg := redisCache.Config{}
		cli, err := redisCache.Init(cfg)

		suite.Require().Error(err)
		suite.Assert().Nil(cli)

	})
	suite.Run("подключение к redis с невалидной конфигурацией", func() {
		cfg := redisCache.Config{
			Addr:     "241421.46334.14241.61253:253532325",
		}
		cli, err := redisCache.Init(cfg)

		suite.Require().Error(err)
		suite.Assert().Nil(cli)

	})
}
