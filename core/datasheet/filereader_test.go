package datasheet_test

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/ff14wed/sibyl/backend/datasheet"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileReader", func() {
	It("properly reads existing files", func() {
		f := new(datasheet.FileReader)
		f.ReadFile("filereader.go", func(r io.Reader) error {
			contents, err := ioutil.ReadAll(r)
			Expect(err).ToNot(HaveOccurred())
			Expect(contents).ToNot(BeEmpty())
			return nil
		})
		Expect(f.Error()).ToNot(HaveOccurred())
	})

	It("errors if the file doesn't exist", func() {
		f := new(datasheet.FileReader)
		f.ReadFile("nonexistent-12345.go", func(r io.Reader) error {
			Fail("Callback shouldn't be called")
			return nil
		})
		f.ReadFile("filereader.go", func(r io.Reader) error {
			Fail("Callback shouldn't be called")
			return nil
		})
		Expect(f.Error()).To(HaveOccurred())
	})

	It("errors if the callback returns an error", func() {
		f := new(datasheet.FileReader)
		f.ReadFile("filereader.go", func(r io.Reader) error {
			return errors.New("Foo")
		})
		f.ReadFile("filereader.go", func(r io.Reader) error {
			Fail("Callback shouldn't be called")
			return nil
		})
		Expect(f.Error()).To(MatchError("callback failure on file filereader.go: Foo"))
	})
})
