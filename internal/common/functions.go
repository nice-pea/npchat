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

func RndPassword() string {
	return gofakeit.Password(true, true, true, true, false, userr.UserPasswordMaxLen)
}

func IsZero[T comparable](v T) bool {
	var zero T
	return v == zero
}

func Zero[T comparable]() T {
	var zero T
	return zero
}
