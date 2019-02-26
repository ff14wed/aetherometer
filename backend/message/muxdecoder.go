package message

import (
	"bytes"
	"io"

	"github.com/ff14wed/xivnet/v3"
)

// DecoderFactory defines a constructor for a FrameDecoder
type DecoderFactory func(r io.Reader) FrameDecoder

// FrameDecoder defines any decoder that can read frames from the input reader
type FrameDecoder interface {
	NextFrame() (*xivnet.Frame, error)
	DiscardDataUntilValid()
}

type decoderStream struct {
	buffer  *bytes.Buffer
	decoder FrameDecoder
}

func newDecoderStream(df DecoderFactory) decoderStream {
	buffer := bytes.NewBuffer(nil)
	return decoderStream{
		buffer:  buffer,
		decoder: df(buffer),
	}
}

// MuxDecoder accepts byte data from any number of sources and stores data
// tagged with a source ID. It decodes full messages from each source and
// provides decoded messages via the NextFrame method.
// This is especially useful for handling multiple sources of data in a single
// goroutine. This same task could be done with goroutines, but it may lead to
// goroutine leaks when it is not clear or impossible to know when these sources
// of data appear or disappear.
// MuxDecoder is not thread-safe. It should either be used in a single-threaded
// context or be handled in a critical section.
type MuxDecoder struct {
	decoderFactory    DecoderFactory
	streams           map[uint32]decoderStream
	lastUpdatedSource uint32
}

// NewMuxDecoder creates a new multiplexer decoder
func NewMuxDecoder(decoderFactory DecoderFactory) *MuxDecoder {
	return &MuxDecoder{
		decoderFactory: decoderFactory,
		streams:        make(map[uint32]decoderStream),
	}
}

// WriteData sends the byte data to the correct decoder stream for handling.
// The next NextFrame call will return the result of calling NextFrame on the
// correct decoder.
// If the decoder stream doesn't exist, it creates a new one.
func (m *MuxDecoder) WriteData(sourceID uint32, message []byte) {
	if _, found := m.streams[sourceID]; !found {
		m.streams[sourceID] = newDecoderStream(m.decoderFactory)
	}
	// Writing to a bytes.Buffer will never error. However, this may panic in very
	// pathologic conditions.
	_, _ = m.streams[sourceID].buffer.Write(message)
	m.lastUpdatedSource = sourceID
}

// NextFrame returns the next valid frame on the last updated decoder. This
// may be called multiple times per WriteData call. If there are no more
// full frames to be read, NextFrame() will return a nil frame with a nil
// error.
func (m *MuxDecoder) NextFrame() (*xivnet.Frame, error) {
	if _, found := m.streams[m.lastUpdatedSource]; !found {
		return nil, nil
	}
	decoder := m.streams[m.lastUpdatedSource].decoder
	for {
		frame, err := decoder.NextFrame()
		if err == nil {
			return frame, nil
		}
		switch err.(type) {
		case xivnet.InvalidHeaderError:
			decoder.DiscardDataUntilValid()
		case xivnet.EOFError:
			return nil, nil
		default:
			return nil, err
		}
	}
}
