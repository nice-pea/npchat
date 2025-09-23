package redisCache

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type JWTIssuanceRegistry struct {
	*redis.Client
}

var (
	ErrEmptySessionID = errors.New("empty session ID")
	ErrEmptyIssueTime = errors.New("empty issue time")
)

func (ir *JWTIssuanceRegistry) RegisterIssueTime(sessionID uuid.UUID, issueTime time.Time) error {
	return nil
}
func (ir *JWTIssuanceRegistry) GetIssueTime(sessionID uuid.UUID) (*time.Time, error) {
	return nil, nil
}
