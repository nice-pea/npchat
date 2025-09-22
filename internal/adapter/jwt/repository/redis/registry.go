package redisCache

import (
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type JWTIssuanceRegistry struct {
	*redis.Client
}

func RegisterIssueTime(sessionID uuid.UUID, issueTime time.Time) error {
	return nil
}
func GetIssueTime(sessionID uuid.UUID) (*time.Time, error) {
	return nil, nil
}
