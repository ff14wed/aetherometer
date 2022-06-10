package hook

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/ff14wed/aetherometer/core/message"
	"github.com/ff14wed/aetherometer/core/stream"
	"github.com/ff14wed/xivnet/v3"
	"github.com/thejerf/suture"
	"go.uber.org/zap"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Stream

// Stream provides the interface for a long running process responsible for
// handling data coming from the hook, as well as sending data to the hook
// when necessary.
type Stream interface {
	suture.Service
	stream.Provider
	fmt.Stringer
}

type hookStream struct {
	streamID uint32
	*suture.Supervisor

	sender      *StreamSender
	frameReader *FrameReader
}

// NewStream creates a new hook Stream
func NewStream(streamID uint32, cfg AdapterConfig, logger *zap.Logger) Stream {
	streamName := fmt.Sprintf("stream-%d", streamID)
	streamLogger := logger.Named(streamName)
	supervisorLogger := streamLogger.Named("supervisor")

	s := &hookStream{
		streamID: streamID,
		Supervisor: suture.New(streamName, suture.Spec{
			Log: func(line string) {
				supervisorLogger.Info(line)
			},
		}),
	}

	hookConn, err := InitializeHook(streamID, cfg)
	if err != nil {
		streamLogger.Error("Failed to initialize hook", zap.Error(err))
		return nil
	}

	pingInterval := 1 * time.Second
	if cfg.HookConfig.PingInterval > 0 {
		pingInterval = time.Duration(cfg.HookConfig.PingInterval)
	}

	ss := NewStreamSender(hookConn, streamLogger)
	sp := NewStreamPinger(ss, pingInterval, streamLogger)
	sr := NewStreamReader(hookConn, streamLogger)
	frameDecoderFactory := func(r io.Reader) message.FrameDecoder {
		if cfg.OodleFactory != nil {
			oodleImpl, err := cfg.OodleFactory.New(streamID)
			if err == nil {
				return xivnet.NewDecoderWithOodle(r, 65535, oodleImpl)
			}
			streamLogger.Error("Failed to initialize oodle, falling back to decoder without Oodle", zap.Error(err))
		}
		return xivnet.NewDecoder(r, 65535)
	}
	fr := NewFrameReader(
		streamID,
		sr.ReceivedEnvelopesListener(),
		frameDecoderFactory,
		streamLogger,
	)

	s.Add(sr)
	s.Add(ss)
	s.Add(sp)
	s.Add(fr)

	s.sender = ss
	s.frameReader = fr

	return s
}

// StreamID returns this stream's ID
func (s *hookStream) StreamID() int {
	return int(s.streamID)
}

// SubscribeIngress provides parsed ingress frames read from the hook
func (s *hookStream) SubscribeIngress() <-chan *xivnet.Frame {
	return s.frameReader.SubscribeIngress()
}

// SubscribeIngress provides parsed egress frames read from the hook
func (s *hookStream) SubscribeEgress() <-chan *xivnet.Frame {
	return s.frameReader.SubscribeEgress()
}

// SendRequest sends a request directly to this stream's hook
func (s *hookStream) SendRequest(req []byte) ([]byte, error) {
	// This particular implementation of SendRequest requires that the request
	// bytes must be directly marshalable to an Envelope
	// Length does not need to be provided
	var env Envelope
	err := json.Unmarshal(req, &env)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshal data to envelope: %s", err)
	}
	s.sender.Send(env.Op, env.Data, env.Additional)
	return []byte(`OK`), nil
}

// InitializeHook injects the specified DLL into the target process and
// attempts to initialize a connection with this hook. It returns the connection
// if successful, and it returns an error if initialization fails.
func InitializeHook(streamID uint32, cfg AdapterConfig) (net.Conn, error) {
	dllPath := cfg.HookConfig.DLLPath
	rpp := cfg.RemoteProcessProvider
	retryInterval := 500 * time.Millisecond
	if cfg.HookConfig.DialRetryInterval > 0 {
		retryInterval = time.Duration(cfg.HookConfig.DialRetryInterval)
	}

	if _, err := os.Stat(dllPath); err != nil {
		return nil, err
	}

	isOwner := true

	err := rpp.InjectDLL(streamID, dllPath)
	if _, ok := err.(DLLAlreadyInjectedError); ok {
		isOwner = false

	} else if err != nil {
		return nil, err
	}

	var conn net.Conn
	pipeName := fmt.Sprintf(`\\.\pipe\xivhook-%d`, streamID)
	dialTimeout := 5 * time.Second

	for i := 0; i < 5; i++ {
		conn, err = rpp.DialPipe(pipeName, &dialTimeout)
		if err == nil {
			return &hookConn{Conn: conn, rpp: rpp, isOwner: isOwner}, nil
		}
		// If we got some sort of error connecting to the pipe, that means
		// the hook hasn't started the pipe server yet. We need to retry in
		// some amount of time.
		time.Sleep(retryInterval)
	}

	return nil, err
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . net.Conn

type hookConn struct {
	net.Conn
	rpp     RemoteProcessProvider
	isOwner bool

	once sync.Once
}

// Write implements the Write interface of a net.Conn. It converts pipe closed
// errors to io.EOF.
func (h *hookConn) Write(p []byte) (int, error) {
	n, err := h.Conn.Write(p)
	if h.rpp.IsPipeClosed(err) {
		return n, io.EOF
	}
	return n, err
}

// Write implements the Read interface of a net.Conn. It converts pipe closed
// errors to io.EOF.
func (h *hookConn) Read(p []byte) (int, error) {
	n, err := h.Conn.Read(p)
	if h.rpp.IsPipeClosed(err) {
		return n, io.EOF
	}
	return n, err
}

// Close implements the Close interface of a net.Conn. It will attempt to send
// an Exit to the connection before closing. It is safe to call Close() more
// than once since subsequent Close() calls will be no-ops.
func (h *hookConn) Close() error {
	var err error
	h.once.Do(func() {
		if h.isOwner {
			// The hook should automatically unload itself from FFXIV after
			// closing the owner core process.
			_, _ = h.Conn.Write(Envelope{Op: OpExit}.Encode())
		}
		err = h.Conn.Close()
	})
	return err
}
