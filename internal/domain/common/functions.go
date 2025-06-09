package common

import (
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

// RndElem возвращает случайный элемент из среза
func RndElem[T any](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	index := rand.Intn(len(slice))
	return slice[index]
}

// RndMapElem возвращает случайный элемент из карты
func RndMapElem[K comparable, V any](m map[K]V) (k K, v V) {
	for k, v := range m {
		return k, v
	}

	return k, v
}

func RndPassword() string {
	return gofakeit.Password(true, true, true, true, false, userr.UserPasswordMaxLen)
}
