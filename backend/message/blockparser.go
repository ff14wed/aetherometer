package message

import (
	"bufio"

	"github.com/ff14wed/xivnet"
)

// FrameDecoder defines any decoder that can read frames from the input reader
type FrameDecoder interface {
	Decode(buf *bufio.Reader) (*xivnet.Frame, error)
	DiscardDataUntilValid(buf *bufio.Reader)
}

// ExtractBlocks reads frames off the reader with the FrameDecoder and
// returns raw blocks that have been extracted from those frames. If there is
// not enough data in the reader to read a full block, it does not consume the
// remaining data in the reader. However if the data in the reader is invalid,
// it discards the bytes in the reader until it is valid.
func ExtractBlocks(buf *bufio.Reader, d FrameDecoder) ([]*xivnet.Block, error) {
	var blocks []*xivnet.Block
	for {
		frame, err := d.Decode(buf)
		switch err {
		case nil:
			blocks = append(blocks, frame.Blocks...)
		case xivnet.ErrInvalidHeader:
			d.DiscardDataUntilValid(buf)
		case xivnet.ErrNotEnoughData:
			return blocks, nil
		default:
			return nil, err
		}
	}
}
