package registerHandler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
)

// loginResultData возвращает результат регистрации входа в виде Map
func loginResultData(session sessionn.Session, user userr.User, jwtIssuer JwtIssuer) fiber.Map {
	data := fiber.Map{
		"session": session,
		"user":    user,
	}

	// Если есть JWT-jwtIssuer, то генерируем JWT
	if jwtIssuer != nil {
		if token, err := jwtIssuer.Issue(session); err != nil {
			slog.Error("Ошибка генерации JWT", "error", err)
		} else {
			data["jwt"] = token
		}
	}

	return data
}
