package redisCache

import (
	"context"
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
	if sessionID == (uuid.UUID{}) {
		return ErrEmptySessionID
	}
	if issueTime == (time.Time{}) {
		return ErrEmptyIssueTime
	}

	status := ir.Client.Set(context.TODO(), sessionID.String(), issueTime, 2*time.Minute)
	_, err := status.Result()
	if err != nil {
		return err
	}

	return nil
}
func (ir *JWTIssuanceRegistry) GetIssueTime(sessionID uuid.UUID) (*time.Time, error) {
	return nil, nil
}
