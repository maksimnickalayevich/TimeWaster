// Package helpers defines all necessary tools
// that helps to manage necessary functionality of TimeWaster.
// GoProc -> Processes management
package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
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
	fileManager     FileManager
	WorkingPath     *string
}

// InitGoProc Initializes all the necessary structs
// and tries to find processes that need time-tracking,
// adds them to struct instance
func (gp *GoProc) InitGoProc() (*[]WasterProcess, error) {
	gp.colorogo = InitColorogo()
	gp.fileManager = FileManager{workingPath: gp.WorkingPath}
	err := gp.fileManager.DetectFiles()
	if err != nil {
		return nil, err
	}

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

	wasterProcesses := gp.FindProcesses(csv)
	if len(*wasterProcesses) == 0 {
		return nil, nil
	}

	return wasterProcesses, nil
}

// FindProcesses finds processes from apps/ folder that are running
func (gp *GoProc) FindProcesses(csvOut string) *[]WasterProcess {
	filesList := gp.fileManager.DetectedFiles

	rawProcessesToTrack := gp.ParseCsv(csvOut, "", *filesList)
	// Create WasterProcesses and add to instance of the app
	gp.createWasterProcess(&rawProcessesToTrack)

	return &gp.processes
}

// ParseCsv parses the csv output of the cmd according to specified separator (sep arg)
func (gp *GoProc) ParseCsv(out string, sep string, processesToTrack map[string]os.FileInfo) []map[string]string {
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

		procName := strings.Replace(splittedLine[0], "\"", "", -1)

		// Don't parse processes that are not necessary for time-tracking
		if _, ok := processesToTrack[procName]; ok == false {
			continue
		}

		processInfo := map[string]string{
			"ProcessName": procName,
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
