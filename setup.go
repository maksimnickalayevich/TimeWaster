package main

import (
	"TimeWaster/helpers"
	"log"
	"os"
)

// Setup Generates apps/ folder to store desktop icons
// or exes of processes to track time
func Setup() (string, error) {
	// Get working dir
	workingPath, err := os.Getwd()
	helpers.HandleError(err, "Unable to get working directory.", true)

	// Create apps/ folder if it doesn't exist
	fullWorkingPath := workingPath + "\\apps"
	exists, err := doesExist(fullWorkingPath)
	helpers.HandleError(err, "Directory apps/ doesn't exist", false)
	if exists == false {
		if err := os.Mkdir(fullWorkingPath, os.ModeDir); err != nil {
			log.Fatalf("Unable to create directory in path %s", workingPath)
			return "", err
		}
	}
	return fullWorkingPath, nil

}

// doesExit checks if entity exits in app installation path
func doesExist(entity string) (bool, error) {
	_, err := os.Stat(entity)
	if os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
