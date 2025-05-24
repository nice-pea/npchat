package http2

// Redirect представляет команду для перенаправления
type Redirect struct {
	URL  string // URL, на который будет выполнено перенаправление
	Code int    // Код статуса HTTP для перенаправления (например, 301, 302)
}
