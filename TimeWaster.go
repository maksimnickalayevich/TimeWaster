package main

import (
	"TimeWaster/helpers"
	"fmt"
	"log"
	"strings"
	"syscall"
	"time"
)

type TimeWaster struct {
	Os            string
	IsRunning     bool
	colorogo      helpers.Colorogo
	processFinder helpers.GoProc
	workingPath   *string
	access        uint32
	timeDelay     uint32
}

func (t *TimeWaster) StartMainLoop(platform string) int {
	colorogo := helpers.InitColorogo()
	t.IsRunning = true
	// Binds colorogo to TimeWaster
	t.colorogo = colorogo

	fmt.Printf(colorogo.Blue+"Welcome to TimeWaster for %s!\n", strings.ToTitle(platform))
	fmt.Println(colorogo.Purple + "What would you like to do?")
	fmt.Println(colorogo.Purple + "1. Track time of your app \n2. Press q/2 to exit")
	for t.IsRunning {
		var choice string
		fmt.Scan(&choice)
		switch choice {
		case "q", "2":
			t.IsRunning = false
		case "1":
			t.TrackTime()
		default:
			fmt.Println(t.colorogo.Red + "You should choose only from suggested options")
		}

	}
	fmt.Print(colorogo.Reset)
	return 1
}

// TrackTime Finds running processes that are set to track time in apps/ folder
func (t *TimeWaster) TrackTime() {
	t.access = syscall.STANDARD_RIGHTS_READ | syscall.PROCESS_QUERY_INFORMATION | syscall.SYNCHRONIZE
	username, err := t.processFinder.GetUsername()
	helpers.HandleError(err, t.colorogo.Red+"Unable to get user", false)
	fmt.Println(t.colorogo.Purple + "Username found: " + username + t.colorogo.Reset)

	fmt.Println(t.colorogo.Yellow + "Please, make sure, that apps you need to track are in apps/ folder" + t.colorogo.Reset)

	// TODO: Create animation of pending app start
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Init processFinder and add working path to it
	t.processFinder = helpers.GoProc{WorkingPath: t.workingPath}

	processes, err := t.processFinder.InitGoProc()
	if processes == nil || err != nil {
		log.Println(t.colorogo.Red + "No process to track were found. Please, be sure your app is running." + t.colorogo.Reset)
		log.Println(t.colorogo.Yellow + "Press any button with option again." + t.colorogo.Reset)
		return
	}

	log.Println(t.colorogo.Yellow + "Starting time-tracking of your apps..." + t.colorogo.Reset)

	t.startConcurrentTrack(time.Duration(2)*time.Second, processes)

	log.Println(t.colorogo.Green + "Your apps were closed." + t.colorogo.Reset)

}

func (t *TimeWaster) startConcurrentTrack(duration time.Duration, apps *[]helpers.WasterProcess) {
	inProcess := make(chan bool)
	go t.checkStatus(duration, inProcess, apps)
	go t.closeRequest(inProcess)
	inProcess <- true
	<-inProcess
	return
}

// TODO: think on what should be returned, assume write to chan time.Time
// checkStatus checks status of apps every duration (seconds), until they are not closed or until user didn't close
// time-tracking
func (t *TimeWaster) checkStatus(duration time.Duration, inProcess chan bool, apps *[]helpers.WasterProcess) {
	running := <-inProcess
	lastCheck := time.Now().Second()
	derefApps := *apps
	if len(derefApps) == 0 {
		close(inProcess)
	}
	log.SetFlags(log.Ltime)
	// Loop to trigger checking status every duration seconds
	// until app is closed and until user didn't stop tracking
	for running && len(derefApps) > 0 {
		// TODO: get not seconds of now, but total time from 1970 in seconds
		now := time.Now().Second()
		difference := now - lastCheck // e.g. 20 - 10 == duration->check; else 20 - 13 == 7 don't check;

		if difference >= int(duration.Seconds()) {
			// TODO: Remove println
			log.Println(t.colorogo.Yellow + "Checking status of running apps" + t.colorogo.Reset)
			for i, app := range derefApps {
				status := t.checkAppStatus(app)
				if !status {
					derefApps = helpers.Remove(derefApps, i)
				}
			}
			log.Println("Updating last exec time...")
			lastCheck = time.Now().Second()
		}
	}
}

func (t *TimeWaster) closeRequest(inProcess chan bool) {
	fmt.Println(t.colorogo.Yellow + "To stop tracking press q or 2" + t.colorogo.Reset)
	var closeEvent string
	_, err := fmt.Scan(&closeEvent)
	if err != nil {
		fmt.Println("Wrong button pressed")
	}
	switch closeEvent {
	case "q", "2":
		inProcess <- false
		return
	default:
		fmt.Println(t.colorogo.Yellow + "Use buttons q or 2 to stop tracking time" + t.colorogo.Reset)
	}
}

// checkAppStatus Checks does app still running, or it's now closed
// Output: running -> true; closed -> false
// N.B. os.FindProcess() doesn't work as it always finds process even if
// it's not running
func (t *TimeWaster) checkAppStatus(app helpers.WasterProcess) bool {
	_, err := syscall.OpenProcess(t.access, false, uint32(app.GetPid()))
	if err != nil {
		return false
	}
	return true
}
