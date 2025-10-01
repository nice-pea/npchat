package redisCache

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Registry struct {
	Cli *redis.Client
	Ttl time.Duration
}

var (
	ErrEmptySessionID = errors.New("empty session ID")
	ErrEmptyIssueTime = errors.New("empty issue time")
)

func (ir *Registry) RegisterIssueTime(sessionID uuid.UUID, issueTime time.Time) error {
	if sessionID == uuid.Nil {
		return ErrEmptySessionID
	}
	if issueTime.IsZero() {
		return ErrEmptyIssueTime
	}

	status := ir.Cli.Set(context.TODO(), sessionID.String(), issueTime, ir.Ttl)
	if _, err := status.Result(); err != nil {
		return err
	}

	return nil
}
func (ir *Registry) IssueTime(sessionID uuid.UUID) (time.Time, error) {
	if sessionID == (uuid.UUID{}) {
		return time.Time{}, ErrEmptySessionID
	}

	var issueTime time.Time

	err := ir.Cli.Get(context.TODO(), sessionID.String()).Scan(&issueTime)
	if errors.Is(err, redis.Nil) {
		return time.Time{}, nil

	}
	return issueTime, nil
}
