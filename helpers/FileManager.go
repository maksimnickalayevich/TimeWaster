package helpers

import (
	"log"
	"os"
	"strings"
)

// FileManager Simple abstraction to detect all the files
// that need time-tracking and writing result of a time track
// to .json file
type FileManager struct {
	DetectedFiles *map[string]os.FileInfo
	workingPath   *string
}

// DetectFiles detect files in the working path
func (fm *FileManager) DetectFiles() error {
	colorogo := InitColorogo()

	dirToCheck := *fm.workingPath
	dirEntry, err := os.ReadDir(dirToCheck)
	if err != nil {
		log.Printf("Error occurred while reading directory %s", dirToCheck)
		return err
	}

	// Collect all file names
	filesNames := make(map[string]os.FileInfo)
	var filesNamesList []string

	for _, d := range dirEntry {
		fileInfo, err := d.Info()
		if err != nil {
			continue
		}
		fileName := fileInfo.Name()
		contains := strings.Contains(fileName, ".lnk")
		if contains {
			fileName = strings.Replace(fileName, ".lnk", "", -1)
		}
		filesNames[fileName] = fileInfo
		filesNamesList = append(filesNamesList, fileName)
	}

	log.Printf(colorogo.Yellow+"Found files: %s"+colorogo.Reset, filesNamesList)

	fm.DetectedFiles = &filesNames

	return nil
}

// GetFileProperties finds properties(detailed info that is used to
// detect process in the whole process list) of the file and returns it
func (fm *FileManager) GetFileProperties(file os.DirEntry) (os.FileInfo, error) {
	fileInfo, err := file.Info()
	if err != nil {
		return nil, err
	}
	return fileInfo, nil
}
