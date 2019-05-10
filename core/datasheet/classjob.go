package datasheet

import (
	"fmt"
	"io"
)

// ClassJobStore stores all of the ClassJob data.
type ClassJobStore map[byte]ClassJob

// ClassJob stores the data for a game ClassJob
type ClassJob struct {
	Key          byte   `datasheet:"key"`
	Name         string `datasheet:"Name"`
	Abbreviation string `datasheet:"Abbreviation"`
}

// PopulateClassJobs will populate the ClassJobStore with ClassJob data provided a
// path to the data sheet for ClassJobs.
func (c *ClassJobStore) PopulateClassJobs(dataReader io.Reader) error {
	*c = make(map[byte]ClassJob)

	var rows []ClassJob
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateClassJobs: %s", err)
	}
	for _, classJob := range rows {
		(*c)[classJob.Key] = classJob
	}
	return nil
}
