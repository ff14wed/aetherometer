package message_test

import (
	"bufio"
	"io"

	"github.com/ff14wed/aetherometer/core/message"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MuxDecoder", func() {
	var (
		muxDecoder        *message.MuxDecoder
		lastCreatedReader *bufio.Reader
	)

	BeforeEach(func() {
		newTestFrameDecoder := func(r io.Reader) message.FrameDecoder {
			reader := bufio.NewReader(r)
			lastCreatedReader = reader
			return &testFrameDecoder{buf: reader}
		}
		muxDecoder = message.NewMuxDecoder(newTestFrameDecoder)
	})

	Context("when there is enough valid data on the buffer", func() {
		BeforeEach(func() {
			muxDecoder.WriteData(1, []byte("PRE-abcPRE-de"))
		})

		It("returns the correct blocks associated with the data", func() {
			f, err := muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(Equal(testFrames["abc"]))
		})

		It("consumes just the valid data on the buffer", func() {
			_, err := muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(lastCreatedReader.Peek(6)).To(Equal([]byte("PRE-de")))
		})

	})

	Context("when there is only invalid data on the buffer", func() {
		BeforeEach(func() {
			muxDecoder.WriteData(1, []byte("invalid-data"))
		})

		It("discards the data until the buffer is smaller than the minimum token size", func() {
			f, err := muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(BeNil())
			Expect(lastCreatedReader.Peek(6)).To(Equal([]byte("d-data")))
		})
	})

	Context("when there is invalid data in between valid frames", func() {
		BeforeEach(func() {
			muxDecoder.WriteData(1, []byte("PRE-abcinvalid-dataPRE-def"))
		})

		It("consumes the invalid data and returns all of the valid frames", func() {
			f, err := muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(Equal(testFrames["abc"]))

			f, err = muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(Equal(testFrames["def"]))

			_, err = lastCreatedReader.Peek(1)
			Expect(err).To(MatchError(io.EOF))
		})
	})

	Context("when there are multiple contiguous frames of data on the buffer", func() {
		BeforeEach(func() {
			muxDecoder.WriteData(1, []byte("PRE-abcPRE-defPRE-abcPRE-g"))
		})

		It("returns all of the correct blocks associated with the data", func() {
			f, err := muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(Equal(testFrames["abc"]))

			f, err = muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(Equal(testFrames["def"]))

			f, err = muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(Equal(testFrames["abc"]))
		})

		It("consumes just the valid data on the buffer", func() {
			_, err := muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			_, err = muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			_, err = muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())

			Expect(lastCreatedReader.Peek(5)).To(Equal([]byte("PRE-g")))
		})
	})

	Context("when data is coming from multiple sources", func() {
		It("correctly separates the data streams when parsing frames", func() {
			By("Writing partial data from sources 1 and 2")
			muxDecoder.WriteData(1, []byte("PRE-ab"))
			f, err := muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(BeNil())

			muxDecoder.WriteData(2, []byte("PRE-de"))
			f, err = muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(BeNil())

			By("Finishing the frame on source 1, should result in only one frame")
			muxDecoder.WriteData(1, []byte("c"))
			f, err = muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(Equal(testFrames["abc"]))

			f, err = muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(BeNil())

			By("Finishing the frame on source 2, should result in only one frame")
			muxDecoder.WriteData(2, []byte("f"))
			f, err = muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(Equal(testFrames["def"]))

			f, err = muxDecoder.NextFrame()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(BeNil())
		})
	})
})
