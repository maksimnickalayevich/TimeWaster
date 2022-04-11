package helpers

import (
	"errors"
	"log"
)

func HandleError(err error, msg string, toPanic bool) {
	if toPanic == false {
		if err != nil {
			log.Printf(msg)
		}
	} else {
		if err != nil {
			log.Panicf(msg)
		}
	}
}

// Remove removes value with index from the slice
func Remove[T any](slice []T, index int) ([]T, error) {
	var newSlice []T
	for i := 0; i < len(slice); i++ {
		if i == index {
			newSlice = append(slice[:i], slice[i+1:]...)
			return newSlice, nil
		}
	}
	return nil, errors.New("index out of range")
}
