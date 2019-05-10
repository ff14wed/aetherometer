package process

import (
	"time"

	"go.uber.org/zap"
)

// Scanner polls the system for all running process names that match a given
// string and emits events to notify of any changes in this list.
type Scanner struct {
	processName string
	scanTicker  <-chan time.Time
	pe          Enumerator
	logger      *zap.Logger

	addProcEventChan chan uint32
	remProcEventChan chan uint32

	procCache map[uint32]struct{}

	stop chan struct{}
}

// NewScanner creates a new process scanner.
// - processName is a substring match query for the process name
// - scanTicker provides a mechanism for the Scanner to iterate on a given
//   interval
// - pe provides the Scanner with an API for enumerating processes
// - eventBufSize determines the size of this event channel
// - logger provides the scanner with a logger
func NewScanner(
	processName string,
	scanTicker <-chan time.Time,
	pe Enumerator,
	eventBufSize int,
	logger *zap.Logger,
) *Scanner {
	return &Scanner{
		processName: processName,
		scanTicker:  scanTicker,
		pe:          pe,
		logger:      logger.Named("scanner"),

		addProcEventChan: make(chan uint32, eventBufSize),
		remProcEventChan: make(chan uint32, eventBufSize),

		procCache: make(map[uint32]struct{}),

		stop: make(chan struct{}),
	}
}

// Serve is responsible for running the process scanner
func (s *Scanner) Serve() {
	s.logger.Info("Running")
	for {
		// Tick first then wait on the ticker
		pidList, err := ListMatchingProcesses(s.processName, s.pe)
		if err != nil {
			s.logger.Error("Nonfatal error", zap.Error(err))
		}
		s.updatePIDs(pidList)
		select {
		case <-s.stop:
			s.logger.Info("Stopping...")
			return
		case <-s.scanTicker:
			continue
		}
	}
}

// Stop closes the process scanner
func (s *Scanner) Stop() {
	close(s.stop)
}

func (s *Scanner) updatePIDs(pidList []uint32) {
	newProcCache := make(map[uint32]struct{})
	for _, pid := range pidList {
		newProcCache[pid] = struct{}{}
		if _, ok := s.procCache[pid]; !ok {
			s.addProcEventChan <- pid
		}
	}
	for pid := range s.procCache {
		if _, ok := newProcCache[pid]; !ok {
			s.remProcEventChan <- pid
		}
	}
	s.procCache = newProcCache
}

// ProcessAddEventListener returns a channel on which subscribers can listen
// for new process events. Scanner does not broadcast events to all subscribers;
// only one will be notified of a given change.
func (s *Scanner) ProcessAddEventListener() <-chan uint32 {
	return s.addProcEventChan
}

// ProcessRemoveEventListener returns a channel on which subscribers can listen
// for close process events. Scanner does not broadcast events to all
// subscribers; only one will be notified of a given change.
func (s *Scanner) ProcessRemoveEventListener() <-chan uint32 {
	return s.remProcEventChan
}
