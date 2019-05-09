package hook

import (
	"io"

	"go.uber.org/zap"
)

//go:generate counterfeiter io.ReadCloser

// StreamReader reads data from the hook connection and decodes it into
// envelopes
type StreamReader struct {
	hookConn io.ReadCloser
	logger   *zap.Logger

	recvChan chan Envelope
	stopDone chan struct{}
}

// NewStreamReader creates a new StreamReader
func NewStreamReader(hookConn io.ReadCloser, logger *zap.Logger) *StreamReader {
	return &StreamReader{
		hookConn: hookConn,
		logger:   logger,

		recvChan: make(chan Envelope),
		stopDone: make(chan struct{}),
	}
}

// Serve runs the service responsible for reading data from the hook connection
// and decoding the data as envelopes.
func (r *StreamReader) Serve() {
	defer close(r.recvChan)
	defer close(r.stopDone)

	logger := r.logger.Named("stream-reader")
	logger.Info("Running")

	d := NewDecoder(r.hookConn, 262144)

	for {
		env, err := d.NextEnvelope()
		if err == nil {
			r.recvChan <- env
		} else if err == io.EOF {
			logger.Info("Stopping...")
			return
		} else {
			logger.Error("reading data from conn", zap.Error(err))
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

// ReceivedEnvelopesListener returns a channel on which consumers can listen
// for envelopes from the hook connection.
func (r *StreamReader) ReceivedEnvelopesListener() <-chan Envelope {
	return r.recvChan
}
