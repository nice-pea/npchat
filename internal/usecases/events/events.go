// Package events определяет интерфейс для работы с событиями.
// Определяет стуктуру для хранения и отправки событий.
//
// Возможно более подходящее названием было бы effects.

package events

// Consumer описывает интерфейс потребителя событий.
type Consumer interface {
	// Consume помещает события в consumer и не может вернуть ошибку
	Consume(events []any)
}

// Buffer представляет структуру для удобного хранения событий
// перед отправкой потребителю
type Buffer struct {
	events []any
}

// Add добавляет событие в буфер
func (ee *Buffer) Add(event any) {
	ee.events = append(ee.events, event)
}

// AddSafety добавляет событие в буфер.
// Если буфер nil, то ничего не делает
func (ee *Buffer) AddSafety(event any) {
	if ee == nil {
		return
	}
	ee.events = append(ee.events, event)
}

// Events возвращает список событий
func (ee *Buffer) Events() []any {
	return ee.events
}
