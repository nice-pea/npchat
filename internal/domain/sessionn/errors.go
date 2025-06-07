package sessionn

import "errors"

var (
	ErrSessionStatusValidate = errors.New("некорректный статус сессии")
	ErrSessionNameEmpty      = errors.New("название сессии не может быть пустым")
)
