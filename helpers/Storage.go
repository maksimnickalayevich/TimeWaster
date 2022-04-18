package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type WasterResult struct {
	OpenTime            time.Time     `json:"-"`
	CloseTime           time.Time     `json:"-"`
	LastExecTime        string        `json:"lastExecTime"`
	StopTime            string        `json:"-"`
	TotalTime           string        `json:"totalTime"`
	LastSession         string        `json:"lastSession"`
	Name                string        `json:"name"`
	UnparsedLastSession time.Duration `json:"-"`
}

func (wr *WasterResult) PopulateResult() {
	wr.LastExecTime = wr.OpenTime.Format(time.RFC822)
	wr.StopTime = wr.CloseTime.Format(time.RFC822)
	wr.LastSession = wr.GetLastSession()
}

func (wr *WasterResult) GetLastSession() string {
	session := wr.CloseTime.Sub(wr.OpenTime)
	wr.UnparsedLastSession = session

	parsed, err := time.ParseDuration(session.String())
	if err != nil {
		log.Fatal("Unable to parse session time")
		return ""
	}

	sessionResult := fmt.Sprintf("%s", parsed.Truncate(time.Second).String())
	return sessionResult
}

// AppendToStorage creates storage based on the type, (default is jsonFile) and writes data to it
func AppendToStorage(storageType StorageType, path string, procs []WasterResult) error {
	colorogo := InitColorogo()
	switch storageType {
	case JsonFile:
		createdFile, err := _createJsonFile(path, &colorogo)
		if err != nil {
			return err
		}
		result := _writeTime(createdFile, procs, &colorogo, true)
		if !result {
			log.Println("Something went wrong while reading/writing the .json file")
		}
	case Database:
		// TODO: create db file + connection to db
		return nil
	}

	return nil
}

func _createJsonFile(path string, colorogo *Colorogo) (*os.File, error) {
	// Check if file already exists
	fullFileName := path + "\\" + "time.json"
	file, err := os.OpenFile(fullFileName, os.O_WRONLY, 0755)
	// if file doesn't exist create it
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(fullFileName)
		if err != nil {
			log.Printf(colorogo.Red+"Unable to create time.json at %s"+colorogo.Reset, path)
			return nil, err
		}
		return file, nil
	}

	return file, err
}

func _writeTime(file *os.File, tracks []WasterResult, colorogo *Colorogo, check bool) bool {
	// Create byte buf to extend it with existing data
	var existingData []byte
	// Check does the file already exist and has some info, or not
	if check {
		content, err := os.ReadFile(file.Name())
		if err != nil {
			log.Println(colorogo.Red + "Unable to read the file time.json" + colorogo.Reset)
			return false
		}
		existingData = append(existingData, content...)
	}

	var processes []WasterResult
	if len(existingData) > 0 {
		err := json.Unmarshal(existingData, &processes)
		if err != nil {
			log.Println("Error happened while unmarshaling json", err)
		}
	}

	updated := updateExisting(processes, tracks)
	marshaledProcesses, err := json.MarshalIndent(updated, "", " ")
	if err != nil {
		log.Println(colorogo.Red + "Unable to marshal data" + colorogo.Reset)
		return false
	}

	// Append data to existing
	_, err = file.Write(marshaledProcesses)
	if err != nil {
		log.Fatal(colorogo.Red + "Unable to write spent time to file" + colorogo.Reset)
		return false
	}

	log.Println(colorogo.Green + "File time.json were successfully updated" + colorogo.Reset)
	err = file.Close()
	if err != nil {
		log.Println(colorogo.Red + "Unable to close the file" + colorogo.Reset)
	}
	return true
}

// updateExisting loops over all the processes (existing and new) and if process
// already wrote to .json file updates it with last info + populates total time field
func updateExisting(existingProcs []WasterResult, newProcs []WasterResult) []WasterResult {
	var updated []WasterResult
	existingMap := make(map[string]WasterResult)
	// Populate map
	for _, ep := range existingProcs {
		existingMap[ep.Name] = ep
	}

	// Add new and update existing
	for _, newEp := range newProcs {
		// Check does it already wrote to file, if yes, update it
		v, ok := existingMap[newEp.Name]
		if ok && !Exists(updated, v.Name) {
			updatedProc := _populateNewProcess(v, newEp)
			updated = append(updated, updatedProc)
			// Add new entry to the updated slice
		} else if !Exists(updated, newEp.Name) {
			updatedProc := _populateNewProcess(newEp, newEp)
			updated = append(updated, updatedProc)
		}
	}

	for name, wr := range existingMap {
		if !Exists(updated, name) {
			updated = append(updated, wr)
		}
	}
	return updated
}

func _populateNewProcess(previous WasterResult, current WasterResult) WasterResult {
	totalTime, _ := time.ParseDuration(previous.TotalTime)
	lastSession := current.UnparsedLastSession
	newTotalTime := fmt.Sprintf("%s", (totalTime + lastSession).Truncate(time.Second).String())

	newProc := WasterResult{
		Name:                current.Name,
		OpenTime:            current.OpenTime,
		CloseTime:           current.CloseTime,
		LastExecTime:        current.LastExecTime,
		StopTime:            current.StopTime,
		TotalTime:           newTotalTime,
		LastSession:         current.LastSession,
		UnparsedLastSession: lastSession,
	}

	return newProc
}
