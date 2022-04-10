// Package helpers defines all necessary tools
// that helps to manage necessary functionality of TimeWaster.
// GoProc -> Processes management
package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"strings"
)

const MIN_PROC_LEN = 5

// GoProc Simple abstraction to handle Processes management in TimeWaster
type GoProc struct {
	username        string
	processes       []WasterProcess
	processesLookup map[string]WasterProcess
	outBuffer       *bytes.Buffer
	errBuffer       *bytes.Buffer
	colorogo        Colorogo
}

// CheckProcesses Tries to find process in list of all processes
// if found -> ptr to initialized GoApp
func (gp *GoProc) CheckProcesses() (*WasterProcess, error) {
	gp.colorogo = InitColorogo()

	// Init buffers for out and err
	gp.outBuffer = &bytes.Buffer{}
	gp.errBuffer = &bytes.Buffer{}

	username, err := gp.GetUsername()
	HandleError(err, gp.colorogo.Red+"Unable to get user", false)
	fmt.Println(gp.colorogo.Purple + "Username found: " + username + gp.colorogo.Reset)

	cmd := exec.Command("tasklist", "/fo", "csv", "/nh")
	cmd.Stdout = gp.outBuffer
	cmd.Stderr = gp.errBuffer
	err = cmd.Run()
	HandleError(err, gp.colorogo.Red+"Unable to find processes", true)

	csv := string(gp.outBuffer.Bytes())

	wasterProcess := gp.IsRunning(csv)
	if wasterProcess == nil {
		return nil, nil
	}

	return wasterProcess, nil
}

// IsRunning Check if all the processes from app/ folder is running
func (gp *GoProc) IsRunning(csvOut string) *WasterProcess {
	// TODO: found process name from file from apps/
	process := "WowClassic.exe"
	processesList := gp.ParseCsv(csvOut, "")
	// Create WasterProcesses and add to instance of the app
	gp.createWasterProcess(&processesList)

	// Get the process if it exists in list of all processes
	wp, ok := gp.processesLookup[process]
	if ok == false {
		return nil
	}

	return &wp
}

// ParseCsv parses the csv output of the cmd according to specified separator (sep arg)
func (gp *GoProc) ParseCsv(out string, sep string) []map[string]string {
	// init default value
	if sep == "" {
		sep = ","
	}
	var processList []map[string]string

	outputLines := strings.Split(out, "\r\n")
	for _, v := range outputLines {
		if v == "" {
			continue
		}
		splittedLine := strings.Split(v, sep)
		if len(splittedLine) < MIN_PROC_LEN {
			log.Printf(gp.colorogo.Red+"Unable to parse process %s"+gp.colorogo.Reset, splittedLine[0])
			continue
		}
		processInfo := map[string]string{
			"ProcessName": strings.Replace(splittedLine[0], "\"", "", -1),
			"Pid":         strings.Replace(splittedLine[1], "\"", "", -1),
			"MemUsage":    strings.Replace(splittedLine[4], "\"", "", -1),
		}
		processList = append(processList, processInfo)
	}

	return processList
}

// createWasterProcess initializes WasterProcess and save the list to GoProc
// in fields GoProc.processes and GoProc.processesLookup
func (gp *GoProc) createWasterProcess(processesListDict *[]map[string]string) {
	if gp.processesLookup == nil {
		gp.processesLookup = make(map[string]WasterProcess)
	}
	for _, p := range *processesListDict {
		// Init new WasterProcess
		wasterProcess := WasterProcess{
			name:     p["ProcessName"],
			pid:      p["Pid"],
			memUsage: p["MemUsage"],
		}
		// Add to list of all processes
		gp.processes = append(gp.processes, wasterProcess)
		// Create fast lookup
		gp.processesLookup[wasterProcess.name] = wasterProcess
	}
}

// GetUsername Gets the username logged in the OS
func (gp *GoProc) GetUsername() (string, error) {
	username, err := user.Current()

	if err != nil {
		return "", err
	}

	unpackedUsername, err := gp.UnpackUsername(username.Username)
	if err != nil {
		return "", err
	}
	gp.username = unpackedUsername

	return unpackedUsername, nil
}

// UnpackUsername Gets the username from full domain e.g. DESKTOP-12345\user
func (gp *GoProc) UnpackUsername(fullDomain string) (string, error) {
	splittedDomain := strings.Split(fullDomain, "\\")
	if len(splittedDomain) < 1 {
		return "", errors.New("unable to get username from OS")
	}

	username := splittedDomain[1]

	return username, nil
}
