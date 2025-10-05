package registerHandler

import (
	"github.com/nice-pea/npchat/internal/domain/sessionn"
)

type JwtIssuer interface {
	Issue(session sessionn.Session) (string, error)
}
