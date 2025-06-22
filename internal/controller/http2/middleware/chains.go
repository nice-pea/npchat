package middleware

import "github.com/nice-pea/npchat/internal/controller/http2"

// ClientAuthChain цепочка обработчиков для клиентских обработчиков с обязательной аутентификацией
var ClientAuthChain = []http2.Middleware{
	RequireRequestID,
	RequireAcceptJson,
	RequireContentTypeJson,
	RequireAuthorizedSession,
}

// ClientPubChain цепочка обработчиков для клиентских обработчиков без аутентификации
var ClientPubChain = []http2.Middleware{
	RequireRequestID,
	RequireAcceptJson,
	RequireContentTypeJson,
}

var EmptyChain []http2.Middleware = nil
