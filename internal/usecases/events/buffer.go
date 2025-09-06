// Package events определяет интерфейс для работы с событиями.
// Определяет стуктуру для хранения и отправки событий.
//
// Возможно более подходящее названием было бы effects.

package events

// Buffer представляет структуру для удобного хранения событий
// перед отправкой потребителю
type Buffer struct {
	events []Event
}

// Add добавляет событие в буфер
func (b *Buffer) Add(event Event) {
	b.events = append(b.events, event)
}

// AddSafety добавляет событие в буфер.
// Если буфер nil, то ничего не делает
func (b *Buffer) AddSafety(event Event) {
	if b == nil {
		return
	}
	b.Add(event)
}

// Events возвращает список событий
func (b *Buffer) Events() []Event {
	return b.events
}
