package utils

import (
	"math/rand/v2"
)

func COALESCE[T comparable](elem ...T) T {
	var empty T
	for _, item := range elem {
		if item != empty {
			return item
		}
	}
	return empty
}

func PickOne[T any](input []T) T {
	return input[rand.IntN(len(input))]
}
