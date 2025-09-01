// Package events определяет интерфейс для работы с событиями.
// Определяет стуктуру для хранения и отправки событий.
//
// Возможно более подходящее названием было бы effects.

package events

type Publisher interface {
	Publish(e *Events) error
}

type Events struct {
	events []any
}

func (ee *Events) Add(event any) {
	ee.events = append(ee.events, event)
}

func (ee *Events) AddSafety(event any) {
	if ee == nil {
		return
	}
	ee.events = append(ee.events, event)
}
