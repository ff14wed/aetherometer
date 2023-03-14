package stream

import (
	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/xivnet/v3"
	"github.com/thejerf/suture"
	"go.uber.org/zap"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Provider

// Provider defines the public facing interface for a provider of a parsed data
// stream. It must provide methods for ingesting data and allow some way of
// controlling the stream.
//
// It is assumed that all blocks produced by this Provider are already parsed
// into the correct xivnet datatype. This is to ensure backwards compatibility
// with older data when the datatype opcodes are updated.
type Provider interface {
	// StreamID returns a unique identifier for the stream. This identifier
	// must be unique across all adapters.
	StreamID() int
	// SubscribeIngress notifies the core of network packets in the ingress
	// direction from this stream.
	SubscribeIngress() <-chan *xivnet.Block
	// SubscribeEgress notifies the core of network packets in the egress
	// direction from this stream.
	SubscribeEgress() <-chan *xivnet.Block
	// SendRequest provides an interface to allow clients to query or control
	// the adapter.
	SendRequest(req []byte) (resp []byte, err error)
}

// Adapter defines an interface that translates data from data sources into
// streams that the core server can consume data from. Each stream provided
// by the adapter is wrapped in a Provider in order for the core server to
// operate on them.
//
// The core server is capable of handling multiple streams of data
// from multiple sources, and they are uniquely identified by stream IDs.
// These stream IDs can correspond to anything to OS processes or just a unique
// identifier for a service or a cluster of services.
//
// Adapters must implement the suture.Service interface so that it can itself be
// started as a long running goroutine. If the adapter must itself parent some
// services, the adapter should itself be or embed a *suture.Supervisor so that
// the non-child nodes of the process tree are only supervisors.
type Adapter interface {
	suture.Service
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . AdapterBuilder

// AdapterBuilder represents a set of methods that instantiates the adapter.
type AdapterBuilder interface {
	// LoadConfig provides the builder with configuration for the adapter. It
	// should return an error if there are any configuration errors other than
	// validation errors. Validation of the config will be handled separately when
	// the configuration is first loaded.
	//
	// If LoadConfig returns an error, Build will not be called for this
	// AdapterBuilder. This error is considered fatal and will cause the server to
	// exit.
	//
	// Adapter authors should add an Enabled field in the adapter configuration
	// struct to conditionally disable the adapter. If the adapter is disabled in
	// the configuration, then this builder is automatically skipped.
	//
	// Adapters are expected to receive their configuration from a section
	// designated by `[adapters.ADAPTER_NAME]`, but it may peek at configuration
	// options outside of this section.
	LoadConfig(config.Config) error

	// Build must return an adapter that is capable of notifying the core
	// with StreamProviders whenever a new stream is to be created. The adapter
	// should also be capable of notifying the core with stream IDs whenever
	// a stream is closed.
	// These StreamProviders should provide the stream's ID and channels
	// that allow the core to consume data from the stream.
	Build(streamUp chan<- Provider, streamDown chan<- int, logger *zap.Logger) Adapter
}

// AdapterInfo lists information about this Adapter and provides a builder
// for the Adapter.
type AdapterInfo struct {
	Name    string
	Builder AdapterBuilder
}
