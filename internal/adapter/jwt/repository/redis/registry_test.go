package redisCache_test

import (
	"time"

	redisCache "github.com/nice-pea/npchat/internal/adapter/jwt/repository/redis"

	"github.com/google/uuid"
)

func (suite *testSuite) Test_JWTIssuanceRegistry() {
	suite.Run("RegisterIssueTime", func() {
		// RegisterIssueTime(sessionID uuid.UUID, issueTime time.Time) error {
		suite.Run("если sessionID пустой то вернется ошибка", func() {
			err := suite.RedisCli.RegisterIssueTime(uuid.UUID{}, time.Now())
			suite.Assert().ErrorAs(err, redisCache.ErrEmptySessionID)
		})
		suite.Run("если issueTime пустой то вернется ошибка", func() {
			err := suite.RedisCli.RegisterIssueTime(uuid.New(), time.Time{})
			suite.Assert().ErrorAs(err, redisCache.ErrEmptyIssueTime)
		})

		suite.Run("новое значение сменит старое", func() {
			sessionId := uuid.New()
			issueTime := time.Now()
			err := suite.RedisCli.RegisterIssueTime(sessionId, issueTime)
			suite.Require().NoError(err)

			issueTimeFromRedis1 := suite.getIssueTime(sessionId)
			suite.Require().Equal(issueTime, issueTimeFromRedis1)

			newIssueTime := time.Now().Add(time.Hour)
			err = suite.RedisCli.RegisterIssueTime(sessionId, newIssueTime)
			suite.Require().NoError(err)

			issueTimeFromRedis2 := suite.getIssueTime(sessionId)
			suite.Require().Equal(newIssueTime, issueTimeFromRedis2)

			suite.Assert().NotEqual(issueTimeFromRedis2, issueTimeFromRedis1)
		})
	})
	suite.Run("GetIssueTime", func() {
		// GetIssueTime(sessionID uuid.UUID) (*time.Time, error)
		suite.Run("из пустого репозитория вернется NULL", func() {

		})

	})
}
