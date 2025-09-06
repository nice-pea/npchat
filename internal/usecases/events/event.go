package events

import (
	"time"

	"github.com/google/uuid"
)

// Event описывает событие
type Event struct {
	Type       string         // Тип события
	CreatedIn  time.Time      // Время создания
	Recipients []uuid.UUID    // Получатели (id пользователей)
	Data       map[string]any // Полезная нагрузка
}
