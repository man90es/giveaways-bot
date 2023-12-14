package utils

import (
	"errors"
	"math/rand"
)

func RandomChoice[T any](items []T) (item T, err error) {
	n := len(items)

	if 1 > n {
		err = errors.New("can't select from an empty array")
	} else {
		item = items[rand.Intn(n)]
	}

	return
}
