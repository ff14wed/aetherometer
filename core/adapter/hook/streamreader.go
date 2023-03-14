package hook

import (
	"io"

	"go.uber.org/zap"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . io.ReadCloser

// StreamReader reads data from the hook connection and decodes it into
// payloads
type StreamReader struct {
	hookConn io.ReadCloser
	logger   *zap.Logger

	recvChan chan Payload
	stopDone chan struct{}
}

// NewStreamReader creates a new StreamReader
func NewStreamReader(hookConn io.ReadCloser, logger *zap.Logger) *StreamReader {
	return &StreamReader{
		hookConn: hookConn,
		logger:   logger.Named("stream-reader"),

		recvChan: make(chan Payload),
		stopDone: make(chan struct{}),
	}
}

// Serve runs the service responsible for reading data from the hook connection
// and decoding the data as payloads.
func (r *StreamReader) Serve() {
	defer close(r.recvChan)
	defer close(r.stopDone)

	r.logger.Info("Running")
	d := NewDecoder(r.hookConn, 262144)

	for {
		env, err := d.NextPayload()
		if err == nil {
			r.recvChan <- env
		} else if err == io.EOF {
			r.logger.Info("Stopping...")
			return
		} else {
			r.logger.Error("reading data from conn", zap.Error(err))
		}
	}
}

// Complete lets the supervisor know it's okay for this process to shut down
// on its own (if the pipe connection shuts down).
func (r *StreamReader) Complete() bool {
	return true
}

// Stop will shutdown this service and wait on it to stop before returning.
func (r *StreamReader) Stop() {
	r.hookConn.Close()
	<-r.stopDone
}

// ReceivedPayloadsListener returns a channel on which consumers can listen
// for payloads from the hook connection.
func (r *StreamReader) ReceivedPayloadsListener() <-chan Payload {
	return r.recvChan
}
