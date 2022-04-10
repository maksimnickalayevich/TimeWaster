package main

import (
	"TimeWaster/helpers"
	"fmt"
	"log"
	"strings"
)

type TimeWaster struct {
	Os            string
	IsRunning     bool
	colorogo      helpers.Colorogo
	processFinder helpers.GoProc
}

func (t *TimeWaster) StartMainLoop(platform string) int {
	colorogo := helpers.InitColorogo()
	t.IsRunning = true
	// Binds colorogo to TimeWaster
	t.colorogo = colorogo

	fmt.Printf(colorogo.Blue+"Welcome to TimeWaster for %s!\n", strings.ToTitle(platform))
	fmt.Println(colorogo.Purple + "What would you like to do?")
	fmt.Println(colorogo.Purple + "1. Track time of your app \n2. See performance statistics of your PC\n3. Press q/3 to exit")
	for t.IsRunning {
		var choice string
		fmt.Scan(&choice)
		switch choice {
		case "q", "3":
			t.IsRunning = false
		case "1":
			t.TrackTime()
		case "2":
			fmt.Println(colorogo.Yellow + "Your choice is 2")
			t.IsRunning = false
		default:
			fmt.Println(t.colorogo.Red + "You should choose only from suggested options")
		}

	}
	fmt.Print(colorogo.Reset)
	return 1
}

// TrackTime Finds running processes that are set to track time if apps/ folder
func (t *TimeWaster) TrackTime() {
	fmt.Println(t.colorogo.Yellow + "Please, make sure, that apps you need to track are in apps/ folder" + t.colorogo.Reset)

	// TODO: Create animation of pending app start
	fmt.Println(t.colorogo.Green + "Waiting for your app to start" + t.colorogo.Reset)
	log.SetFlags(log.Ldate | log.Lshortfile)

	appProcess, err := t.processFinder.CheckProcesses()
	if appProcess == nil {
		log.Printf(t.colorogo.Red + "No process to track were found. Please, be sure your app is running." + t.colorogo.Reset)
		log.Printf(t.colorogo.Yellow + "Press any button with option again.")
		return
	}
	helpers.HandleError(err, "Error occurred", true)

	fmt.Printf("Found proccess: %+v", appProcess)
}
