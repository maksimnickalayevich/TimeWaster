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
