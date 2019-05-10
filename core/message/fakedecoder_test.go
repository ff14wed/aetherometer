package message_test

import (
	"bufio"
	"errors"
	"time"

	"github.com/ff14wed/xivnet/v3"
)

var testFrames = map[string]*xivnet.Frame{
	"abc": &xivnet.Frame{
		Time: time.Unix(12, 0),
		Blocks: []*xivnet.Block{
			&xivnet.Block{
				Length:    123,
				SubjectID: 1234, CurrentID: 5678,
				IPCHeader: xivnet.IPCHeader{Opcode: 0x90},
				Data:      xivnet.GenericBlockDataFromBytes([]byte{1, 2, 3, 4}),
			},
			&xivnet.Block{
				Length:    456,
				SubjectID: 5678, CurrentID: 5678,
				IPCHeader: xivnet.IPCHeader{Opcode: 0x91},
				Data:      xivnet.GenericBlockDataFromBytes([]byte{5, 6, 7, 8}),
			},
		},
	},
	"def": &xivnet.Frame{
		Blocks: []*xivnet.Block{
			&xivnet.Block{
				Length:    789,
				SubjectID: 2345, CurrentID: 5678,
				IPCHeader: xivnet.IPCHeader{Opcode: 0x92},
				Data:      xivnet.GenericBlockDataFromBytes([]byte{9, 0, 1, 2}),
			},
		},
	},
}

type testFrameDecoder struct {
	buf *bufio.Reader
}

func (d *testFrameDecoder) NextFrame() (*xivnet.Frame, error) {
	token, err := d.buf.Peek(7)
	if err != nil {
		return nil, xivnet.EOFError{}
	}
	if string(token[0:4]) != "PRE-" {
		return nil, xivnet.InvalidHeaderError{}
	}
	key := token[4:7]
	_, _ = d.buf.Discard(7)
	if f, ok := testFrames[string(key)]; ok {
		return f, nil
	}
	return nil, errors.New("invalid data")
}

func (d *testFrameDecoder) DiscardDataUntilValid() {
	for {
		token, err := d.buf.Peek(7)
		if err != nil {
			return
		}
		if string(token[0:4]) == "PRE-" {
			return
		}

		_, _ = d.buf.Discard(1)
	}
}
