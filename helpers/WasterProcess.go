package helpers

import (
	"log"
	"strconv"
)

type WasterProcess struct {
	name     string
	pid      string
	memUsage string
}

func (wp *WasterProcess) GetPid() int {
	parsedPid, err := strconv.ParseInt(wp.pid, 10, 32)
	if err != nil {
		log.Printf("Unable to parse app pid %s - %s", wp.GetName(), wp.pid)
	}
	
	return int(parsedPid)
}

func (wp *WasterProcess) GetName() string {
	return wp.name
}
