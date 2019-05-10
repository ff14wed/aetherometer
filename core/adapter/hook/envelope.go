package hook

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

// Constants for the ops field in an Envelope
const (
	OpDebug = iota
	OpPing
	OpExit
	OpRecv
	OpSend
)

// Envelope defines the message format used to communicate with the hook
// Length is the length of the entire envelope, including the length field.
// The size of the Envelope is 9 bytes + len(Additional).
type Envelope struct {
	Length     uint32
	Op         byte
	Data       uint32
	Additional []byte
}

// DecodeEnvelope transforms byte data to an Envelope
func DecodeEnvelope(data []byte) Envelope {
	e := Envelope{}
	e.Length = binary.LittleEndian.Uint32(data[0:4])
	e.Op = data[4]
	e.Data = binary.LittleEndian.Uint32(data[5:9])
	e.Additional = make([]byte, len(data[9:]))
	copy(e.Additional, data[9:])
	return e
}

// EncodeEnvelope transforms an Envelope to byte data
func (e Envelope) Encode() []byte {
	buf := make([]byte, len(e.Additional)+9)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(len(buf)))
	buf[4] = e.Op
	binary.LittleEndian.PutUint32(buf[5:9], e.Data)
	copy(buf[9:], e.Additional)
	return buf
}

// ErrInvalidLength is returned whenever data is corrupted in the byte
// stream and we have potentially faulty data.
var ErrInvalidLength = errors.New("Invalid length encountered in byte stream")

// Decoder is responsible for reading bytes from the provided reader
// and decoding the data into envelopes
type Decoder struct {
	reader *bufio.Reader
}

// NewDecoder creates a new decoder instance given an io.Reader and a buffer
// size. Data read from the Reader will be buffered to store partial data
// as it comes in.
func NewDecoder(r io.Reader, bufSize int) *Decoder {
	return &Decoder{
		reader: bufio.NewReaderSize(r, bufSize),
	}
}

func (d *Decoder) consumeBytes(numBytes uint32) {
	_, _ = d.reader.Discard(int(numBytes))
}

// NextEnvelope consumes an some amount of data on the Reader and returns
// the next full Envelope. Since the Envelope isn't by itself a robust format
// of data transmission, the decoder might or might not recover from decoding
// faulty data. Since the intended io.Reader is a named pipe connection, the
// chance of a failure happening is really low.
//
// If the length is at least readable, but it's too small, it will discard
// the faulty data and continue attempting to read Envelopes. However,
// there is no recovery path if the data is corrupted in other ways, and
// subsequent calls to NextEnvelope will return the same thing.
func (d *Decoder) NextEnvelope() (Envelope, error) {
	lengthBytes, err := d.reader.Peek(4)
	if err != nil {
		return Envelope{}, err
	}

	length := binary.LittleEndian.Uint32(lengthBytes)

	if length < 4 {
		length = 4
	}
	if length < 9 {
		d.consumeBytes(length)
		return Envelope{}, ErrInvalidLength
	}

	envBytes, err := d.reader.Peek(int(length))
	if err != nil {
		return Envelope{}, err
	}

	env := DecodeEnvelope(envBytes)
	d.consumeBytes(length)
	return env, nil
}
