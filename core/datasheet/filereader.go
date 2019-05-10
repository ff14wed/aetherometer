package datasheet

import (
	"fmt"
	"io"
	"os"
)

// FileReader handles reading files and errors if they don't exist
type FileReader struct {
	lastError error
}

type fileReaderCallback func(io.Reader) error

// ReadFile reads the file provided and calls the callback function with the
// io.Reader. It won't do anything if a previous ReadFile has encountered
// an error.
func (f *FileReader) ReadFile(fileName string, cb fileReaderCallback) {
	if f.lastError != nil {
		return
	}
	file, err := os.Open(fileName)
	if err != nil {
		f.lastError = fmt.Errorf("reading file %s: %s", fileName, err)
		return
	}
	defer file.Close()
	err = cb(file)
	if err != nil {
		f.lastError = fmt.Errorf("callback failure on file %s: %s", fileName, err)
		return
	}
}

func (f *FileReader) Error() error {
	return f.lastError
}
