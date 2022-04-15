package helpers

import (
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
func Remove[T any](slice []T, index int) []T {
	var newSlice []T
	for i := 0; i < len(slice); i++ {
		if i == index {
			newSlice = append(slice[:i], slice[i+1:]...)
			return newSlice
		}
	}
	return nil
}
