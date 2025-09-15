package http2

import (
	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/controller/http2/register_handler"
)

type RequiredJwt interface {
	register_handler.JwtIssuer
	middleware.JwtParser
}
