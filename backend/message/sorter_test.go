package message_test

import (
	"io/ioutil"

	"github.com/ff14wed/sibyl/backend/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sorter", func() {
	Describe("Sort", func() {
		It("creates and returns the buffered reader that received the new message", func() {
			sorter := message.NewSorter()
			reader := sorter.Sort(1, []byte("Hello"))
			Expect(reader).ToNot(BeNil())
			Expect(ioutil.ReadAll(reader)).To(Equal([]byte("Hello")))
		})

		It("correctly sorts each message into their respective mailboxes", func() {
			sorter := message.NewSorter()
			reader := sorter.Sort(1, []byte("Hello"))
			Expect(reader).ToNot(BeNil())
			reader2 := sorter.Sort(2, []byte("Foo"))
			Expect(reader2).ToNot(BeNil())
			reader1 := sorter.Sort(1, []byte(" World"))
			Expect(reader1).To(Equal(reader))
			Expect(ioutil.ReadAll(reader)).To(Equal([]byte("Hello World")))
			Expect(ioutil.ReadAll(reader2)).To(Equal([]byte("Foo")))
		})
	})
})
