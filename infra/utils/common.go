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
	if len(input) == 0 {
		panic("PickOne empty input")
	}
	if len(input) == 1 {
		return input[0]
	}
	return input[rand.IntN(len(input))]
}
