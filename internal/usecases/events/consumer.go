package events

// Consumer описывает интерфейс потребителя событий.
type Consumer interface {
	// Consume помещает события в consumer и не может вернуть ошибку
	Consume(events []Event)
}
