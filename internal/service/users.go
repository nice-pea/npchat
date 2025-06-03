package service

import (
	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

// Users сервис, объединяющий случаи использования(юзкейсы) в контексте агрегата пользователей
type Users struct {
	Providers    OAuthProviders // Карта провайдеров OAuth
	Repo         userr.Repository
	SessionsRepo domain.SessionsRepository // Репозиторий сессий пользователей
}
