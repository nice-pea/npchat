package registerHandler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/nice-pea/npchat/internal/common"
)

func Info(router *fiber.App, buildInfo common.BuildInfo) {
	router.Get(
		"/info",
		func(ctx *fiber.Ctx) error {
			return ctx.JSON(fiber.Map{
				"version":    buildInfo.Version,
				"build_date": buildInfo.BuildDate,
				"commit":     buildInfo.Commit,
			})
		})
}
