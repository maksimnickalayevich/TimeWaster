package main

import (
	"TimeWaster/helpers"
	"context"
	"fmt"
	"log"
	"syscall"
	"time"
)

type TimeWaster struct {
	Os            string
	IsRunning     bool
	colorogo      helpers.Colorogo
	processFinder helpers.GoProc
	workingPath   *string
	results       []helpers.WasterResult
	closeTimes    map[string]time.Time
}

func (t *TimeWaster) StartMainLoop(platform string) int {
	colorogo := helpers.InitColorogo()
	t.IsRunning = true
	t.closeTimes = make(map[string]time.Time)

	// Binds colorogo to TimeWaster
	t.colorogo = colorogo

	fmt.Printf(colorogo.Blue+"Welcome to TimeWaster for %s!\n", platform)
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

			log.Println(t.colorogo.Green + "Updating .json with your spent time" + t.colorogo.Reset)

			err := helpers.AppendToStorage(helpers.JsonFile, *t.workingPath, t.results)
			if err != nil {
				log.Fatal(t.colorogo.Red + "Error happened while updating time.json file" + t.colorogo.Reset)
				return 0
			}

		default:
			fmt.Println(t.colorogo.Red + "You should choose only from suggested options: 1; q/2")
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

	log.SetFlags(log.Ldate | log.Ltime)

	// Init processFinder and add working path to it
	t.processFinder = helpers.GoProc{WorkingPath: t.workingPath}

	processes, err := t.processFinder.InitGoProc()
	if processes == nil || err != nil {
		log.Println(t.colorogo.Red + "No process to track were found. Please, be sure your app is running." + t.colorogo.Reset)
		log.Println(t.colorogo.Yellow + "Press any button with option again." + t.colorogo.Reset)
		return
	}
	deref := *processes

	log.Println(t.colorogo.Yellow + "Starting time-tracking of your apps..." + t.colorogo.Reset)

	var results []helpers.WasterResult
	results = t.Convert(deref)

	t.startConcurrentTrack(time.Duration(5)*time.Second, processes)

	// Update with populated values
	results = t.updateResults(&results)

	t.results = results
}

// Convert Converts WasterProcesses to WasterResults and initializes them with basic data.
// e.g. {Name, OpenTime}
func (t *TimeWaster) Convert(what []helpers.WasterProcess) []helpers.WasterResult {
	var converted []helpers.WasterResult

	for _, wp := range what {
		newWasterRes := helpers.WasterResult{
			Name:     wp.GetName(),
			OpenTime: time.Now(),
		}
		converted = append(converted, newWasterRes)
	}

	return converted
}

// updateResults updates results of a time tracking (CloseTime + other fields of WasterResult struct)
func (t *TimeWaster) updateResults(res *[]helpers.WasterResult) []helpers.WasterResult {
	var newRes []helpers.WasterResult

	for _, proc := range *res {
		closeTime, _ := t.closeTimes[proc.Name]
		proc.CloseTime = closeTime
		proc.PopulateResult()
		newRes = append(newRes, proc)
	}

	return newRes
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
				for i, app := range derefApps {
					// Save in times to check what
					t.closeTimes[app.GetName()] = time.Now()
					status := t.checkAppStatus(app)
					if !status {
						derefApps = helpers.Remove(derefApps, i)
						// Update close time by proc name
						t.closeTimes[app.GetName()] = time.Now()
					}
				}
				lastCheck = time.Now().UTC().Unix()
			}
		}
	}
	log.Println(t.colorogo.Blue + "Your apps were closed" + t.colorogo.Reset)
}

// closeRequest runs the loop with listening the user input. If input == q/2 closes goroutine with time tracking via context
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
func (t *TimeWaster) checkAppStatus(app helpers.WasterProcess) bool {
	p, _ := syscall.OpenProcess(helpers.ACCESS, false, uint32(app.GetPid()))
	var exitCode uint32
	exitCodeErr := syscall.GetExitCodeProcess(p, &exitCode)
	if exitCodeErr != nil {
		log.Printf(t.colorogo.Red+"Smth went wrong %e"+t.colorogo.Reset, exitCodeErr)
		return false
	}
	if exitCode != helpers.StillAlive {
		return false
	}
	return true
}
