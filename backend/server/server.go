package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ff14wed/sibyl/backend/config"
	"go.uber.org/zap"
)

// Server handles serving the backend API for Sibyl
type Server struct {
	logger *zap.Logger

	s          *http.Server
	ctx        context.Context
	cancelFunc context.CancelFunc
	address    net.Addr
	ready      chan struct{}

	serveMux *http.ServeMux
}

// New initializes a new instance of the backend server
func New(
	cfg config.Config,
	logger *zap.Logger,
) *Server {
	serveMux := http.NewServeMux()

	s := &Server{
		logger: logger.Named("http-server"),

		s: &http.Server{
			Addr:    fmt.Sprintf("localhost:%d", cfg.APIPort),
			Handler: serveMux,
		},
		ready:    make(chan struct{}),
		serveMux: serveMux,
	}
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())

	return s
}

// AddHandler adds a handler to this server
func (s *Server) AddHandler(path string, handler http.Handler) {
	s.serveMux.Handle(path, handler)
}

// sets a TCP keep-alive timeout on connections so that dead
// TCP connections eventually go away
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	_ = tc.SetKeepAlive(true)
	_ = tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// Serve is responsible for running the backend server
func (s *Server) Serve() {
	ln, err := net.Listen("tcp", s.s.Addr)
	if err != nil {
		s.logger.Error("Error starting listener", zap.Error(err))
		return
	}
	s.address = ln.Addr()

	s.logger.Info("Running", zap.String("address", s.Address().String()))
	close(s.ready)

	err = s.s.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
	if err != nil && err != http.ErrServerClosed {
		s.logger.Error("Error serving", zap.Error(err))
		return
	}
}

// WaitUntilReady blocks until the server has started listening
func (s *Server) WaitUntilReady() {
	<-s.ready
}

// Stop will shutdown the backend server, and will timeout within 1 second
func (s *Server) Stop() {
	s.logger.Info("Stopping...")
	s.cancelFunc()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := s.s.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Error shutting down", zap.Error(err))
		return
	}
}

// Address returns the address of this server. This value is only valid if
// the server has been started, otherwise it'll be a nil Addr interface value.
func (s *Server) Address() *net.TCPAddr {
	if s.address == nil {
		return nil
	}
	return s.address.(*net.TCPAddr)
}
