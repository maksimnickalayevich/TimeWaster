package main

import (
	"TimeWaster/helpers"
	"fmt"
	"runtime"
)

func main() {
	colors := helpers.InitColorogo()

	fmt.Println(colors.Yellow + "Setting up apps/ folder.")
	_, err := Setup()
	helpers.HandleError(err, "Something went wrong while creating apps/ folder.", true)
	fmt.Println(colors.Green + "Directory apps/ created.")

	timeWaster := TimeWaster{}

	timeWaster.StartMainLoop(runtime.GOOS)
}
