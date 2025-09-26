package registerHandler

import (
	"time"

	"github.com/google/uuid"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
)

type JwtIssuer interface {
	Issue(session sessionn.Session) (string, error)
}

type JWTIssuanceRegistry interface {
	RegisterIssueTime(sessionID uuid.UUID, issueTime time.Time) error
	GetIssueTime(sessionID uuid.UUID) (*time.Time, error)
}
