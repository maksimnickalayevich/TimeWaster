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

// Remove removes value from the slice; Not generic solution
func Remove(slice []WasterProcess, uniqueValue interface{}) []WasterProcess {
	var newSlice []WasterProcess

	switch uniqueValue.(type) {
	// Used to remove by unique string value (ObjectId, hash etc.)
	case string:
		index := 0
		for index < len(slice) {
			if slice[index].name != uniqueValue {
				newSlice = append(newSlice, slice[index])
			}
			index++
		}
	// Used to remove by index
	case int:
		newSlice = append(slice[:uniqueValue.(int)], slice[uniqueValue.(int)+1:]...)
	}

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
