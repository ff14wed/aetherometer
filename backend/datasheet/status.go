package datasheet

import (
	"fmt"
	"io"
)

// StatusStore stores all of the Status data.
type StatusStore map[uint32]Status

// Status stores the data for a game Status
type Status struct {
	Key         uint32 `datasheet:"key"`
	Name        string `datasheet:"Name"`
	Description string `datasheet:"Description"`
}

// PopulateStatuses will populate the StatusStore with Status data provided a
// path to the data sheet for Statuses.
func (s *StatusStore) PopulateStatuses(dataReader io.Reader) error {
	*s = make(map[uint32]Status)

	var rows []Status
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateStatuses: %s", err)
	}
	for _, status := range rows {
		(*s)[status.Key] = status
	}
	return nil
}
