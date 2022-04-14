package main

import (
	"TimeWaster/helpers"
	"runtime"
)

func main() {
	setupDir, err := Setup()
	helpers.HandleError(err, "Something went wrong while creating apps/ folder.", true)

	timeWaster := TimeWaster{workingPath: &setupDir}

	timeWaster.StartMainLoop(runtime.GOOS)
}
