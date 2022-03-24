package process

import (
	"fmt"
	"strings"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Enumerator

// Enumerator defines the interface for enumerating processes on the
// system
type Enumerator interface {
	EnumerateProcesses() (map[uint32]string, error)
}

// ListMatchingProcesses lists all IDs for processes that match the given string
func ListMatchingProcesses(match string, pe Enumerator) ([]uint32, error) {
	procMap, err := pe.EnumerateProcesses()
	if err != nil {
		return nil, fmt.Errorf("EnumerateProcesses error: %s", err.Error())
	}

	var pids []uint32
	for pid, procName := range procMap {
		if strings.Contains(strings.ToLower(procName), strings.ToLower(match)) {
			pids = append(pids, pid)
		}
	}
	return pids, nil
}
