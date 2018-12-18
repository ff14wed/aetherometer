package message_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/ff14wed/sibyl/backend/message"

	"github.com/ff14wed/xivnet"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testFrames = map[string]*xivnet.Frame{
	"abc": &xivnet.Frame{
		Blocks: []*xivnet.Block{
			&xivnet.Block{
				Length: 123,
				Header: xivnet.BlockHeader{
					SubjectID: 1234,
					CurrentID: 5678,
					Opcode:    0x90,
				},
				Data: []byte{1, 2, 3, 4},
			},
			&xivnet.Block{
				Length: 456,
				Header: xivnet.BlockHeader{
					SubjectID: 5678,
					CurrentID: 5678,
					Opcode:    0x91,
				},
				Data: []byte{5, 6, 7, 8},
			},
		},
	},
	"def": &xivnet.Frame{
		Blocks: []*xivnet.Block{
			&xivnet.Block{
				Length: 789,
				Header: xivnet.BlockHeader{
					SubjectID: 2345,
					CurrentID: 5678,
					Opcode:    0x92,
				},
				Data: []byte{9, 0, 1, 2},
			},
		},
	},
}

type testFrameDecoder struct{}

func (d *testFrameDecoder) Decode(buf *bufio.Reader) (*xivnet.Frame, error) {
	token, err := buf.Peek(7)
	if err != nil {
		return nil, xivnet.ErrNotEnoughData
	}
	if string(token[0:4]) != "PRE-" {
		return nil, xivnet.ErrInvalidHeader
	}
	key := string(token[4:7])
	_, _ = buf.Discard(7)
	if f, ok := testFrames[key]; ok {
		return f, nil
	}
	return nil, errors.New("invalid data")
}

func (d *testFrameDecoder) DiscardDataUntilValid(buf *bufio.Reader) {
	for {
		token, err := buf.Peek(7)
		if err != nil {
			return
		}
		if string(token[0:4]) == "PRE-" {
			return
		}

		_, _ = buf.Discard(1)
	}
}

var _ = Describe("Block Parser", func() {
	Describe("ExtractBlocks", func() {
		var (
			buf    *bytes.Buffer
			reader *bufio.Reader
		)

		BeforeEach(func() {
			buf = bytes.NewBuffer(nil)
			reader = bufio.NewReader(buf)
		})

		Context("when there is enough valid data on the buffer", func() {
			BeforeEach(func() {
				buf.WriteString("PRE-abcPRE-de")
			})

			It("returns the correct blocks associated with the data", func() {
				blocks, err := message.ExtractBlocks(reader, new(testFrameDecoder))
				Expect(err).ToNot(HaveOccurred())
				Expect(blocks).To(Equal(testFrames["abc"].Blocks))
			})

			It("consumes just the valid data on the buffer", func() {
				_, err := message.ExtractBlocks(reader, new(testFrameDecoder))
				Expect(err).ToNot(HaveOccurred())
				Expect(reader.Peek(6)).To(Equal([]byte("PRE-de")))
			})
		})

		Context("when there is only invalid data on the buffer", func() {
			BeforeEach(func() {
				buf.WriteString("invalid-data")
			})

			It("discards the data until the buffer is smaller than the minimum token size", func() {
				blocks, err := message.ExtractBlocks(reader, new(testFrameDecoder))
				Expect(err).ToNot(HaveOccurred())
				Expect(blocks).To(BeEmpty())
				d, err := reader.Peek(6)
				Expect(err).ToNot(HaveOccurred())
				Expect(d).To(Equal([]byte("d-data")))
			})
		})

		Context("when there is invalid data in between valid blocks", func() {
			var expectedBlocks []*xivnet.Block

			BeforeEach(func() {
				buf.WriteString("PRE-abcinvalid-dataPRE-def")
				expectedBlockSlices := [][]*xivnet.Block{
					testFrames["abc"].Blocks,
					testFrames["def"].Blocks,
				}
				expectedBlocks = nil
				for _, s := range expectedBlockSlices {
					expectedBlocks = append(expectedBlocks, s...)
				}
			})

			It("consumes the invalid data and returns all of the valid blocks", func() {
				blocks, err := message.ExtractBlocks(reader, new(testFrameDecoder))
				Expect(err).ToNot(HaveOccurred())
				Expect(blocks).To(Equal(expectedBlocks))
				_, err = reader.Peek(1)
				Expect(err).To(MatchError(io.EOF))
			})
		})

		Context("when there are multiple contiguous blocks of data on the buffer", func() {
			var expectedBlocks []*xivnet.Block

			BeforeEach(func() {
				buf.WriteString("PRE-abcPRE-defPRE-abcPRE-g")
				expectedBlockSlices := [][]*xivnet.Block{
					testFrames["abc"].Blocks,
					testFrames["def"].Blocks,
					testFrames["abc"].Blocks,
				}
				expectedBlocks = nil
				for _, s := range expectedBlockSlices {
					expectedBlocks = append(expectedBlocks, s...)
				}
			})

			It("returns the correct blocks associated with the data", func() {
				blocks, err := message.ExtractBlocks(reader, new(testFrameDecoder))
				Expect(err).ToNot(HaveOccurred())
				Expect(blocks).To(Equal(expectedBlocks))
			})

			It("consumes just the valid data on the buffer", func() {
				_, err := message.ExtractBlocks(reader, new(testFrameDecoder))
				Expect(err).ToNot(HaveOccurred())
				Expect(reader.Peek(5)).To(Equal([]byte("PRE-g")))
			})
		})
	})
})
