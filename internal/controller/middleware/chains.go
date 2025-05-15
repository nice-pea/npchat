package middleware

import "github.com/saime-0/nice-pea-chat/internal/controller/http2"

// ClientAuthChain цепочка обработчиков для клиентских обработчиков с обязательной аутентификацией
var ClientAuthChain = []http2.Middleware{
	RequireRequestID,
	RequireAcceptJson,
	RequireAuthorizedSession,
}

// ClientPubChain цепочка обработчиков для клиентских обработчиков без аутентификации
var ClientPubChain = []http2.Middleware{
	RequireRequestID,
	RequireAcceptJson,
}

var EmptyChain []http2.Middleware = nil
