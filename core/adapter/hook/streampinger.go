package hook

import (
	"time"

	"go.uber.org/zap"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . HookDataSender

// HookDataSender defines the interface that allows the sending of data
// to the hook connection from the adapter.
type HookDataSender interface {
	Send(op byte, data uint32, additional []byte)
}

// StreamPinger sends a ping through the hook connection to make sure it's still
// alive
type StreamPinger struct {
	hds          HookDataSender
	pingInterval time.Duration
	logger       *zap.Logger

	stop     chan struct{}
	stopDone chan struct{}
}

// NewStreamPinger returns a new StreamPinger
func NewStreamPinger(hds HookDataSender, pingInterval time.Duration, logger *zap.Logger) *StreamPinger {
	return &StreamPinger{
		hds:          hds,
		pingInterval: pingInterval,
		logger:       logger.Named("stream-pinger"),

		stop:     make(chan struct{}),
		stopDone: make(chan struct{}),
	}
}

// Serve runs the service responsible for the periodic pinging of the hook
// connection.
func (p *StreamPinger) Serve() {
	defer close(p.stopDone)
	p.logger.Info("Running")
	t := time.NewTicker(p.pingInterval)

	for {
		select {
		case <-t.C:
			p.hds.Send(OpPing, 0, nil)
		case <-p.stop:
			p.logger.Info("Stopping...")
			return
		}
	}
}

// Stop will shutdown this service and wait on it to stop before returning.
func (p *StreamPinger) Stop() {
	close(p.stop)
	<-p.stopDone
}
