package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
)

func UserID(ctx *fiber.Ctx) uuid.UUID {
	return ctx.Locals(middleware.CtxKeyUserID).(uuid.UUID)
}
func SessionID(ctx *fiber.Ctx) uuid.UUID {
	return ctx.Locals(middleware.CtxKeySessionID).(uuid.UUID)
}

// ParamsUUID возвращает значение из пути запроса как uuid
func ParamsUUID(ctx *fiber.Ctx, name string) uuid.UUID {
	val, _ := uuid.Parse(ctx.Params(name))
	return val
}
