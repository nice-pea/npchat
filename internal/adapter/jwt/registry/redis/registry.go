package redisRegistry

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Registry репозиторий для хранения и получения временных меток выпуска токенов
type Registry struct {
	Client *redis.Client // Клиент Redis
	Ttl    time.Duration // Время жизни записей в кэше
}

// Ошибки модуля
var (
	ErrEmptySessionID = errors.New("SessionID отсутствует")      // Пустой идентификатор сессии
	ErrEmptyIssueTime = errors.New("время создания отсутствует") // Нулевая временная метка
)

// RegisterIssueTime регистрирует время анулирования токена в Redis
func (ir *Registry) RegisterIssueTime(sessionID uuid.UUID, issueTime time.Time) error {
	// Проверка валидности входных данных
	if sessionID == uuid.Nil {
		return ErrEmptySessionID
	}
	if issueTime.IsZero() {
		return ErrEmptyIssueTime
	}

	// Установка значения в Redis с TTL
	status := ir.Client.Set(context.Background(), sessionID.String(), issueTime, ir.Ttl)
	if _, err := status.Result(); err != nil {
		return err // Возвращает ошибку Redis при сбое
	}

	return nil // Успешное выполнение
}

// IssueTime возвращает время последнего анулирования токена из Redis
func (ir *Registry) IssueTime(sessionID uuid.UUID) (time.Time, error) {
	// Проверка валидности идентификатора
	if sessionID == uuid.Nil {
		return time.Time{}, ErrEmptySessionID
	}

	var issueTime time.Time

	// Получение значения из Redis
	err := ir.Client.Get(context.Background(), sessionID.String()).Scan(&issueTime)

	// Обработка случаев отсутствия ключа
	if errors.Is(err, redis.Nil) {
		return time.Time{}, nil // Ключ не найден, но это не ошибка
	}

	if err != nil {
		return time.Time{}, err // Ошибка Redis
	}

	return issueTime, nil // Успешное получение значения
}
