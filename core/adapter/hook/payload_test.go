package hook_test

import (
	"bytes"
	"io"

	"github.com/ff14wed/aetherometer/core/adapter/hook"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var testPayload = hook.Payload{
	Length:  0,
	Op:      hook.OpPing,
	Channel: 123456789,
	Data:    []byte{1, 2, 3, 4},
}

var expectedPayload = testPayload

func init() {
	expectedPayload.Length = 13
}

var expectedPayloadBytes = []byte{
	13, 0, 0, 0, // Correct length is 13
	1,
	0x15, 0xCD, 0x5B, 0x07,
	1, 2, 3, 4,
}

var _ = Describe("Payload", func() {
	Describe("Encode", func() {
		It("encodes the payload data to a byte stream with the correct length field", func() {
			Expect(testPayload.Encode()).To(Equal(expectedPayloadBytes))
		})
	})

	Describe("DecodePayload", func() {
		It("decodes the byte stream to payload data", func() {
			Expect(hook.DecodePayload(expectedPayloadBytes)).To(Equal(expectedPayload))
		})

		It("panics when given fewer than 9 bytes", func() {
			testBytes := []byte{1, 2, 3, 4, 5, 6, 7, 8}
			Expect(func() {
				hook.DecodePayload(testBytes)
			}).To(Panic())
		})
	})

	Describe("Decoder", func() {
		It("consumes and decodes the next full payload on the reader", func() {
			buf := bytes.NewBuffer(append(expectedPayloadBytes, 13, 14, 15, 16))
			d := hook.NewDecoder(buf, 1024)
			env, err := d.NextPayload()
			Expect(err).ToNot(HaveOccurred())
			Expect(env).To(Equal(expectedPayload))
		})

		It("stores partial data in its buffer until the data is complete", func() {
			buf := bytes.NewBuffer(append(expectedPayloadBytes, 13, 0, 0, 0))
			d := hook.NewDecoder(buf, 1024)
			_, err := d.NextPayload()
			Expect(err).ToNot(HaveOccurred())
			_, err = buf.Write(expectedPayloadBytes[4:])
			Expect(err).ToNot(HaveOccurred())

			env, err := d.NextPayload()
			Expect(err).ToNot(HaveOccurred())
			Expect(env).To(Equal(expectedPayload))
		})

		It("errors and does not consume bytes if the reader doesn't enough data to read the length", func() {
			buf := bytes.NewBuffer([]byte{14, 0})
			d := hook.NewDecoder(buf, 1024)
			_, err := d.NextPayload()
			Expect(err).To(MatchError(io.EOF))

			_, err = d.NextPayload()
			Expect(err).To(MatchError(io.EOF))
		})

		It("errors and does not consume bytes if there is not enough data in the reader to read all {length} bytes", func() {
			buf := bytes.NewBuffer(append([]byte{14}, expectedPayloadBytes[1:]...))
			d := hook.NewDecoder(buf, 1024)
			_, err := d.NextPayload()
			Expect(err).To(MatchError(io.EOF))

			_, err = d.NextPayload()
			Expect(err).To(MatchError(io.EOF))
		})

		It("errors if the length is smaller than 4, but still consumes 4 bytes", func() {
			buf := bytes.NewBuffer(append([]byte{0, 0, 0, 0}, expectedPayloadBytes...))
			d := hook.NewDecoder(buf, 1024)
			_, err := d.NextPayload()
			Expect(err).To(MatchError(hook.ErrInvalidLength))

			env, err := d.NextPayload()
			Expect(err).ToNot(HaveOccurred())
			Expect(env).To(Equal(expectedPayload))
		})

		It("errors if the length is 4, but still consumes 4 bytes", func() {
			buf := bytes.NewBuffer(append([]byte{4, 0, 0, 0}, expectedPayloadBytes...))
			d := hook.NewDecoder(buf, 1024)
			_, err := d.NextPayload()
			Expect(err).To(MatchError(hook.ErrInvalidLength))

			env, err := d.NextPayload()
			Expect(err).ToNot(HaveOccurred())
			Expect(env).To(Equal(expectedPayload))
		})

		It("errors if the length is less than 9, but still consumes that number of bytes", func() {
			buf := bytes.NewBuffer(append([]byte{8, 0, 0, 0, 1, 2, 3, 4}, expectedPayloadBytes...))
			d := hook.NewDecoder(buf, 1024)
			_, err := d.NextPayload()
			Expect(err).To(MatchError(hook.ErrInvalidLength))

			env, err := d.NextPayload()
			Expect(err).ToNot(HaveOccurred())
			Expect(env).To(Equal(expectedPayload))
		})
	})
})
