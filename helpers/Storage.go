package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type WasterResult struct {
	OpenTime     time.Time `json:"-"`
	CloseTime    time.Time `json:"-"`
	LastExecTime string    `json:"lastExecTime"`
	StopTime     string    `json:"stopTime"`
	TotalTime    string    `json:"totalTime"`
	LastSession  string    `json:"lastSession"`
	Name         string    `json:"name"`
}

func (wr *WasterResult) PopulateResult() {
	wr.LastExecTime = wr.OpenTime.Format(time.RFC822)
	wr.StopTime = wr.CloseTime.Format(time.RFC822)
	wr.LastSession = wr.GetLastSession()
}

func (wr *WasterResult) GetLastSession() string {
	session := wr.CloseTime.Sub(wr.OpenTime)
	parsed, err := time.ParseDuration(session.String())
	if err != nil {
		log.Fatal("Unable to parse session time")
		return ""
	}

	sessionResult := fmt.Sprintf("Session: %v", parsed.Truncate(time.Second).String())
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
	file, err := os.Create(path + "/time.json")
	if err != nil {
		log.Printf(colorogo.Red+"Unable to create time.json at %s"+colorogo.Reset, path)
		return nil, err
	}
	return file, err
}

func _writeTime(file *os.File, procs []WasterResult, colorogo *Colorogo, check bool) bool {
	// Create byte buf to extend it with existing data
	var existingData []byte

	// Check does the file already exist and has some info, or not
	if check {
		// Read file and if not empty update info
		outer, err := os.ReadFile(file.Name())
		if err != nil {
			log.Fatal(colorogo.Red + "Unable to read the file!" + colorogo.Reset)
			return false
		}
		// TODO: file is not properly read, check this one
		log.Printf("File .json: %b\n", outer)
		existingData = outer
	}

	// Write to a file again
	// TODO: Take care of adding totalTime field
	var dataToWrite []byte

	for _, proc := range procs {
		jsonifiedWr, err := json.Marshal(proc)
		if err != nil {
			log.Fatal(colorogo.Red + "Unable to convert process result to json" + colorogo.Reset)
		}
		dataToWrite = append(existingData, jsonifiedWr...)
	}

	// Append data to existing
	err := os.WriteFile(file.Name(), dataToWrite, os.FileMode(os.O_WRONLY))
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
