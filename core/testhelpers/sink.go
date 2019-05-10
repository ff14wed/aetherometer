package testhelpers

import (
	"github.com/onsi/gomega/gbytes"
	"go.uber.org/zap"
)

var _ zap.Sink = new(LogBuffer)

// LogBuffer wraps a gbytes.Buffer to match the zap.Sink interface
type LogBuffer struct {
	internal *gbytes.Buffer
}

// Sync is an empty implementation of the sync method
func (LogBuffer) Sync() error {
	return nil
}

// Reset resets the internal gbytes Buffer to a clean state
func (l *LogBuffer) Reset() {
	l.internal = gbytes.NewBuffer()
}

// Buffer returns the internal buffer
func (l *LogBuffer) Buffer() *gbytes.Buffer {
	return l.internal
}

// Write implements the io.Writer interface
func (l *LogBuffer) Write(p []byte) (n int, err error) {
	return l.internal.Write(p)
}

// Close implements the io.Closer interface
func (l *LogBuffer) Close() error {
	return l.internal.Close()
}
