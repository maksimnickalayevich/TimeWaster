package main

import (
	"TimeWaster/helpers"
	"fmt"
	"log"
	"strings"
	"syscall"
)

type TimeWaster struct {
	Os            string
	IsRunning     bool
	colorogo      helpers.Colorogo
	processFinder helpers.GoProc
	workingPath   *string
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

// TrackTime Finds running processes that are set to track time in apps/ folder
func (t *TimeWaster) TrackTime() {
	username, err := t.processFinder.GetUsername()
	helpers.HandleError(err, t.colorogo.Red+"Unable to get user", false)
	fmt.Println(t.colorogo.Purple + "Username found: " + username + t.colorogo.Reset)

	fmt.Println(t.colorogo.Yellow + "Please, make sure, that apps you need to track are in apps/ folder" + t.colorogo.Reset)

	// TODO: Create animation of pending app start
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Init processFinder and add working path to it
	t.processFinder = helpers.GoProc{WorkingPath: t.workingPath}

	processes, err := t.processFinder.InitGoProc()
	if processes == nil {
		log.Println(t.colorogo.Red + "No process to track were found. Please, be sure your app is running." + t.colorogo.Reset)
		log.Println(t.colorogo.Yellow + "Press any button with option again." + t.colorogo.Reset)
		return
	}
	helpers.HandleError(err, "Error occurred", true)

	log.Println(t.colorogo.Yellow + "Starting time-tracking of your apps..." + t.colorogo.Reset)
	derefProcesses := *processes
	// TODO: check status periodically
	for len(derefProcesses) > 0 {
		for i, app := range derefProcesses {
			status := t.checkAppStatus(app)
			if status == false {
				derefProcesses, err = helpers.Remove(derefProcesses, i)
			}
		}
	}

	log.Println(t.colorogo.Green + "Your apps were closed." + t.colorogo.Reset)

}

// checkAppStatus Checks does app still running, or it's now closed
// Output: running -> true; closed -> false
// N.B. os.FindProcess() doesn't work as it always finds process even if
// it's not running
func (t *TimeWaster) checkAppStatus(app helpers.WasterProcess) bool {
	const da = syscall.STANDARD_RIGHTS_READ | syscall.PROCESS_QUERY_INFORMATION | syscall.SYNCHRONIZE
	_, err := syscall.OpenProcess(da, false, uint32(app.GetPid()))
	if err != nil {
		return false
	}
	return true
}
