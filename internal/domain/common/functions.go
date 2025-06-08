package common

import "math/rand"

// RndElem возвращает случайный элемент из среза
func RndElem[T any](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	index := rand.Intn(len(slice))
	return slice[index]
}
