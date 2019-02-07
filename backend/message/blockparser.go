package message

import (
	"bufio"

	"github.com/ff14wed/xivnet/v2"
	"github.com/ff14wed/xivnet/v2/datatypes"
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
		if err == nil {
			blocks = append(blocks, frame.Blocks...)
			continue
		}
		switch err.(type) {
		case xivnet.InvalidHeaderError:
			d.DiscardDataUntilValid(buf)
		case xivnet.EOFError:
			return blocks, nil
		default:
			return nil, err
		}
	}
}

// DedupMyMovementBlocks returns a filtered list that does not contain
// duplicated outgoing movement blocks.
// This handler assumes that the list of blocks resulted from a single
// pass of the decoder, and therefore the blocks are too temporally
// close to be useful for us. This is especially an issue since the FFXIV
// client spams the server with movement blocks during casts to ensure
// any movement will interrupt casts.
// This processing must happen here since as soon as we handle blocks
// individually it's too late to dedup.
func DedupMyMovementBlocks(blocks []*xivnet.Block) []*xivnet.Block {
	var prevBlock *xivnet.Block
	deduping := false
	dedupedBlocks := make([]*xivnet.Block, 0, len(blocks))

	for _, b := range blocks {
		switch b.Header.Opcode {
		case datatypes.MyMovementOpcode, datatypes.MyMovement2Opcode:
			deduping = true
		default:
			if deduping {
				dedupedBlocks = append(dedupedBlocks, prevBlock)
				deduping = false
			}
			dedupedBlocks = append(dedupedBlocks, b)
		}
		prevBlock = b
	}
	if deduping {
		dedupedBlocks = append(dedupedBlocks, prevBlock)
	}
	return dedupedBlocks
}
