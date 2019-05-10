package store

import "time"

type providerConfig struct {
	queryTimeout      time.Duration
	updateBufferSize  int
	eventBufferSize   int
	requestBufferSize int
}

// Option defines an optional configuration parameter to the constructor of the
// Provider
type Option func(p *providerConfig)

// WithQueryTimeout sets the timeout for read accesses to the store provider.
// Since all store acceses happen on the same goroutine, it is possible for
// misbehaving updates to block reads. The query methods on the provider will
// timeout after this duration if this scenario happens.
//
// The default value is 5 seconds.
func WithQueryTimeout(t time.Duration) Option {
	return func(p *providerConfig) {
		p.queryTimeout = t
	}
}

// WithUpdateBufferSize sets the size of the store provider channel that is
// responsible for receiving updates. There shouldn't be any real reason to
// change this unless the provider is very slow at consuming updates.
//
// The default value is 10000.
func WithUpdateBufferSize(size int) Option {
	return func(p *providerConfig) {
		p.updateBufferSize = size
	}
}

// WithEventBufferSize sets the size of the outgoing event buffer. There
// shouldn't be any real reason to change this unless the event consumer is very
// slow at consuming updates.
//
// The default value is 10000.
func WithEventBufferSize(size int) Option {
	return func(p *providerConfig) {
		p.eventBufferSize = size
	}
}

// WithRequestBufferSize sets the size of the internal request queue.
// be any real reason to change this unless the the core API is getting hammered
// with requests.
//
// The default value is 10.
func WithRequestBufferSize(size int) Option {
	return func(p *providerConfig) {
		p.requestBufferSize = size
	}
}
