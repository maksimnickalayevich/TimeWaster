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
	newSlice := append(slice[:index], slice[index+1:]...)
	return newSlice
}

// Exists checks if value already exist in the slice
// useful for structs with unique value of type string
func Exists(slice []WasterResult, value string) bool {
	for _, wr := range slice {
		if value == wr.Name {
			return true
		}
	}
	return false
}

// FilterEmpty Filters given slice from empty values with shifting
func FilterEmpty(slice []WasterResult) {
	// TODO: Filter empty value for this slice will be struct{}

}
