package message_test

import (
	"time"

	"github.com/ff14wed/aetherometer/core/message"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProcessBlocks", func() {
	setupBlocksList := func(cb func([]*xivnet.Block) []*xivnet.Block) []*xivnet.Block {
		blocks := []*xivnet.Block{
			&xivnet.Block{
				Length:    123,
				IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressClientTriggerOpcode},
				Data:      xivnet.GenericBlockDataFromBytes([]byte{1, 2, 3, 4}),
			},
			&xivnet.Block{
				Length:    123,
				IPCHeader: xivnet.IPCHeader{Opcode: datatypes.CastingOpcode},
				Data:      xivnet.GenericBlockDataFromBytes([]byte{1, 2, 3, 4}),
			},
		}
		blocks = cb(blocks)
		return append(blocks, &xivnet.Block{
			Length:    123,
			IPCHeader: xivnet.IPCHeader{Opcode: datatypes.CastingOpcode},
			Data:      xivnet.GenericBlockDataFromBytes([]byte{5, 6, 7, 8}),
		})
	}

	It("sets the frame's timestamp on all blocks", func() {
		f := &xivnet.Frame{
			Time: time.Unix(12, 0),
			Blocks: []*xivnet.Block{
				&xivnet.Block{
					Length:    123,
					IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressMovementOpcode},
					Data:      xivnet.GenericBlockDataFromBytes([]byte{1, 2, 3, 4}),
				},
				&xivnet.Block{
					Length: 123,
					Data:   xivnet.GenericBlockDataFromBytes([]byte{5, 6, 7, 8}),
				},
			},
		}
		expectedBlocks := []*xivnet.Block{
			&xivnet.Block{
				Length: 123,
				IPCHeader: xivnet.IPCHeader{
					Opcode: datatypes.EgressMovementOpcode,
					Time:   time.Unix(12, 0),
				},
				Data: xivnet.GenericBlockDataFromBytes([]byte{1, 2, 3, 4}),
			},
			&xivnet.Block{
				Length:    123,
				IPCHeader: xivnet.IPCHeader{Time: time.Unix(12, 0)},
				Data:      xivnet.GenericBlockDataFromBytes([]byte{5, 6, 7, 8}),
			},
		}

		blocks := message.ProcessBlocks(f)
		Expect(blocks).To(Equal(expectedBlocks))
	})

	It("removes duplicate EgressMovement blocks, leaves the last EgressMovement block and unrelated blocks untouched", func() {
		origBlocks := setupBlocksList(func(blocks []*xivnet.Block) []*xivnet.Block {
			for i := 0; i < 20; i++ {
				blocks = append(blocks, &xivnet.Block{
					Length:    123,
					IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressMovementOpcode},
					Data:      xivnet.GenericBlockDataFromBytes([]byte{1, 2, 3, 4}),
				})
			}
			return append(blocks, &xivnet.Block{
				Length:    123,
				IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressMovementOpcode},
				Data:      xivnet.GenericBlockDataFromBytes([]byte{5, 6, 7, 8}),
			})
		})
		expectedBlocks := setupBlocksList(func(blocks []*xivnet.Block) []*xivnet.Block {
			return append(blocks, &xivnet.Block{
				Length:    123,
				IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressMovementOpcode},
				Data:      xivnet.GenericBlockDataFromBytes([]byte{5, 6, 7, 8}),
			})
		})

		f := &xivnet.Frame{Blocks: origBlocks}
		dedupedBlocks := message.ProcessBlocks(f)
		Expect(dedupedBlocks).To(Equal(expectedBlocks))
	})

	It("removes duplicate EgressInstanceMovement blocks, leaves the last EgressInstanceMovement block and unrelated blocks untouched", func() {
		origBlocks := setupBlocksList(func(blocks []*xivnet.Block) []*xivnet.Block {
			for i := 0; i < 20; i++ {
				blocks = append(blocks, &xivnet.Block{
					Length:    123,
					IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressInstanceMovementOpcode},
					Data:      xivnet.GenericBlockDataFromBytes([]byte{1, 2, 3, 4}),
				})
			}
			blocks = append(blocks, &xivnet.Block{
				Length:    123,
				IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressInstanceMovementOpcode},
				Data:      xivnet.GenericBlockDataFromBytes([]byte{5, 6, 7, 8}),
			})
			return blocks
		})
		expectedBlocks := setupBlocksList(func(blocks []*xivnet.Block) []*xivnet.Block {
			blocks = append(blocks, &xivnet.Block{
				Length:    123,
				IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressInstanceMovementOpcode},
				Data:      xivnet.GenericBlockDataFromBytes([]byte{5, 6, 7, 8}),
			})
			return blocks
		})

		f := &xivnet.Frame{Blocks: origBlocks}
		dedupedBlocks := message.ProcessBlocks(f)
		Expect(dedupedBlocks).To(Equal(expectedBlocks))
	})

	Context("when there are only EgressMovement blocks", func() {
		It("returns only the last block", func() {
			var origBlocks []*xivnet.Block
			for i := 0; i < 20; i++ {
				origBlocks = append(origBlocks, &xivnet.Block{
					Length:    123,
					IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressMovementOpcode},
					Data:      xivnet.GenericBlockDataFromBytes([]byte{1, 2, 3, 4}),
				})
			}
			origBlocks = append(origBlocks, &xivnet.Block{
				Length:    123,
				IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressMovementOpcode},
				Data:      xivnet.GenericBlockDataFromBytes([]byte{5, 6, 7, 8}),
			})
			expectedBlocks := []*xivnet.Block{
				&xivnet.Block{
					Length:    123,
					IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressMovementOpcode},
					Data:      xivnet.GenericBlockDataFromBytes([]byte{5, 6, 7, 8}),
				},
			}
			f := &xivnet.Frame{Blocks: origBlocks}
			dedupedBlocks := message.ProcessBlocks(f)
			Expect(dedupedBlocks).To(Equal(expectedBlocks))
		})
	})

	Context("when there is one EgressMovement block", func() {
		It("just leaves the entire block list alone", func() {
			origBlocks := setupBlocksList(func(blocks []*xivnet.Block) []*xivnet.Block {
				return append(blocks, &xivnet.Block{
					Length:    123,
					IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressMovementOpcode},
					Data:      xivnet.GenericBlockDataFromBytes([]byte{5, 6, 7, 8}),
				})
			})
			f := &xivnet.Frame{Blocks: origBlocks}
			dedupedBlocks := message.ProcessBlocks(f)
			Expect(dedupedBlocks).To(Equal(origBlocks))
		})
	})

	Context("when there are no EgressMovement or EgressInstanceMovement blocks", func() {
		It("just leaves the entire block list alone", func() {
			origBlocks := setupBlocksList(func(b []*xivnet.Block) []*xivnet.Block {
				return b
			})
			f := &xivnet.Frame{Blocks: origBlocks}
			dedupedBlocks := message.ProcessBlocks(f)
			Expect(dedupedBlocks).To(Equal(origBlocks))
		})
	})
})
