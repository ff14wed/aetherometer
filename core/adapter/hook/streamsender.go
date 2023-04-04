package hook

import (
	"errors"
	"io"
	"sync/atomic"

	"go.uber.org/zap"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . io.WriteCloser

// StreamSender listens for requests and sends data to the hook.
type StreamSender struct {
	hookConn io.WriteCloser
	logger   *zap.Logger

	isClosed uint32
	sendChan chan Payload
	stopDone chan struct{}
}

// NewStreamSender creates a new StreamSender
func NewStreamSender(hookConn io.WriteCloser, logger *zap.Logger) *StreamSender {
	return &StreamSender{
		hookConn: hookConn,
		logger:   logger.Named("stream-sender"),

		sendChan: make(chan Payload),
		stopDone: make(chan struct{}),
	}
}

// Serve runs the service responsible for handling requests to send data to the
// hook connection.
func (s *StreamSender) Serve() {
	defer close(s.stopDone)
	s.logger.Info("Running")
	for {
		payload, ok := <-s.sendChan
		if !ok {
			s.logger.Info("Stopping...")
			return
		}
		envBytes := payload.Encode()
		n, err := s.hookConn.Write(envBytes)
		if err != nil {
			s.logger.Error("writing bytes to conn", zap.Error(err))
		}
		if n < len(envBytes) {
			err := errors.New("wrote less than the payload length")
			s.logger.Error("writing bytes to conn", zap.Error(err))
		}
	}
}

// Stop will shutdown this service and wait on it to stop before returning.
func (s *StreamSender) Stop() {
	atomic.StoreUint32(&s.isClosed, 1)
	close(s.sendChan)
	<-s.stopDone
}

// Send queues a request to send the data as an Payload through the hook
// connection.
func (s *StreamSender) Send(op byte, channel uint32, data []byte) {
	if atomic.LoadUint32(&s.isClosed) == 0 {
		s.sendChan <- Payload{Op: op, Channel: channel, Data: data}
	}
}
