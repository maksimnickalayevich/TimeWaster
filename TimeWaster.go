package main

import (
	"TimeWaster/helpers"
	"context"
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
	// TODO: Move this to separate struct in order to have access to the menu from every part of a programm
	fmt.Println(colorogo.Purple + "What would you like to do?")
	fmt.Println(colorogo.Purple + "1. Track time of your app \n2. Press q/2 to exit")
	for t.IsRunning {
		var choice string
		_, err := fmt.Scan(&choice)
		if err != nil {
			log.Println(t.colorogo.Red + "Invalid input, unable to parse input" + t.colorogo.Reset)
		}
		switch choice {
		case "q", "2":
			t.IsRunning = false
		case "1":
			t.TrackTime()
		default:
			fmt.Println(t.colorogo.Red + "You should choose only from suggested options: 1; q/2")
		}
	}
	fmt.Print(colorogo.Reset)

	log.Println(t.colorogo.Green + "Updating .json with your spent time" + t.colorogo.Reset)
	return 1
}

// TrackTime Finds running processes that are set to track time in apps/ folder
func (t *TimeWaster) TrackTime() {
	username, err := t.processFinder.GetUsername()

	helpers.HandleError(err, t.colorogo.Red+"Unable to get user", false)

	fmt.Println(t.colorogo.Purple + "Username found: " + username + t.colorogo.Reset)

	fmt.Println(t.colorogo.Yellow + "Please, make sure, that apps you need to track are in apps/ folder" + t.colorogo.Reset)

	log.SetFlags(log.Ldate | log.Ltime)

	// Init processFinder and add working path to it
	t.processFinder = helpers.GoProc{WorkingPath: t.workingPath}

	processes, err := t.processFinder.InitGoProc()
	if processes == nil || err != nil {
		log.Println(t.colorogo.Red + "No process to track were found. Please, be sure your app is running." + t.colorogo.Reset)
		log.Println(t.colorogo.Yellow + "Press any button with option again." + t.colorogo.Reset)
		return
	}

	log.Println(t.colorogo.Yellow + "Starting time-tracking of your apps..." + t.colorogo.Reset)

	t.startConcurrentTrack(time.Duration(5)*time.Second, processes)
}

func (t *TimeWaster) startConcurrentTrack(duration time.Duration, apps *[]helpers.WasterProcess) {
	inProcess := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	// Create goroutines
	go t.checkStatus(duration, inProcess, apps, ctx)
	go t.closeRequest(inProcess, cancel)
	inProcess <- true
	// block until write in close request loop
	<-inProcess

	return
}

// TODO: think on what should be returned, assume write to chan time.Time
// checkStatus checks status of apps every duration (seconds), until they are not closed or until user didn't close
// time-tracking
func (t *TimeWaster) checkStatus(duration time.Duration, inProcess chan bool, apps *[]helpers.WasterProcess, ctx context.Context) {
	running := <-inProcess
	lastCheck := time.Now().UTC().Unix()
	derefApps := *apps

	// Loop to trigger checking status every duration seconds
	// until app is closed and until user didn't stop tracking
	for running && len(derefApps) > 0 {
		select {
		case <-ctx.Done():
			log.Println(t.colorogo.Green + "You triggered full stop of time-tracking" + t.colorogo.Reset)
			return
		default:
			now := time.Now().UTC().Unix()
			difference := now - lastCheck

			if difference >= int64(duration.Seconds()) {
				log.Println(t.colorogo.Yellow + "Checking status of running apps" + t.colorogo.Reset)
				for i, app := range derefApps {
					status := t.checkAppStatus(app)
					if !status {
						derefApps = helpers.Remove(derefApps, i)
					}
				}
				lastCheck = time.Now().UTC().Unix()
			}
		}
	}
	log.Println(t.colorogo.Blue + "Your apps were closed" + t.colorogo.Reset)
}

func (t *TimeWaster) closeRequest(inProcess chan bool, ctxCancel context.CancelFunc) {
	log.Println(t.colorogo.Yellow + "To stop tracking press q or 2" + t.colorogo.Reset)
	for {
		res := t.catchKeyboardEvent(inProcess, ctxCancel)
		if !res {
			break
		}
	}
}

// catchKeyboardEvent if false -> close the check status goroutine; if true -> continue asking the input
func (t *TimeWaster) catchKeyboardEvent(inProcess chan bool, ctxCancel context.CancelFunc) bool {
	// Catch keyboard button press
	var closeEvent string
	_, err := fmt.Scan(&closeEvent)
	if err != nil {
		fmt.Println("Wrong button pressed")
	}
	switch closeEvent {
	case "q", "2":
		inProcess <- false
		close(inProcess) // Free resources
		ctxCancel()      // Close goroutine with time-tracking
		return false
	default:
		fmt.Println(t.colorogo.Yellow + "Use buttons q or 2 to stop tracking time" + t.colorogo.Reset)
		return true
	}
}

// checkAppStatus Checks does app still running, or it's now closed
// Output: running -> true; closed -> false
// N.B. os.FindProcess() doesn't work as it always finds process even if
// it's not running
func (t *TimeWaster) checkAppStatus(app helpers.WasterProcess) bool {
	// TODO: This still doesn't define was it closed or not
	const da = syscall.STANDARD_RIGHTS_READ | syscall.PROCESS_QUERY_INFORMATION | syscall.SYNCHRONIZE
	p, err := syscall.OpenProcess(da, false, uint32(app.GetPid()))
	if err != nil {
		return false
	}
	fmt.Println(p)
	return true
}
