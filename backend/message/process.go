package message

import (
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

// ProcessBlocks returns a processed list of blocks from the provided frame.
//
// Note this will not parse the blocks, as parsing could return errors. These
// errors should be handled on a per-block basis.
//
// The processing will apply a higher resolution timestamp to each block and
// filter out excessive egress movement blocks.
// Because the FFXIV client spams the server with many temporally-close movement
// blocks during casts (presumably to ensure movement will interrupt casts), not
// many of these blocks are very useful for our purposes. Therefore, we filter
// them out during processing.
//
// This processing must happen here since as soon as we handle blocks
// individually it's too late to dedup.
func ProcessBlocks(f *xivnet.Frame) []*xivnet.Block {
	f.CorrectTimestamps(f.Time)
	var prevBlock *xivnet.Block
	deduping := false
	dedupedBlocks := make([]*xivnet.Block, 0, len(f.Blocks))

	for _, b := range f.Blocks {
		switch b.Opcode {
		case datatypes.EgressMovementOpcode, datatypes.EgressInstanceMovementOpcode:
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
