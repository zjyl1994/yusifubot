package utils

import (
	"encoding/json"
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

func MarshalToJsonNoError(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
