package models_test

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/ff14wed/sibyl/backend/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Custom Types", func() {
	Describe("Stream", func() {
		var stream *models.Stream

		BeforeEach(func() {
			stream = &models.Stream{
				PID: 1234,
				EntitiesMap: map[uint64]*models.Entity{
					1: &models.Entity{ID: 1, Name: "FooBar", Index: 2},
					2: &models.Entity{ID: 2, Name: "Baah", Index: 1},
					3: nil,
				},
			}
		})

		Describe("Entities", func() {
			It("returns all entities found on the stream, sorted in order by index", func() {
				Expect(stream.Entities()).To(Equal([]models.Entity{
					*stream.EntitiesMap[2],
					*stream.EntitiesMap[1],
				}))
			})
		})
	})

	Describe("Timestamp", func() {
		It("marshals the provided time to the time since the Unix epoch in milliseconds", func() {
			t := time.Unix(101, 302000000)
			m := models.MarshalTimestamp(t)
			b := new(bytes.Buffer)
			m.MarshalGQL(b)
			Expect(b.String()).To(Equal("101302"))
		})
	})

	Describe("Uint", func() {
		It("marshals the provided uint to string", func() {
			m := models.MarshalUint(123456789000000)
			b := new(bytes.Buffer)
			m.MarshalGQL(b)
			Expect(b.String()).To(Equal("123456789000000"))
		})

		It("unmarshals the string to the expected uint", func() {
			u, err := models.UnmarshalUint("123456789000000")
			Expect(err).ToNot(HaveOccurred())
			Expect(u).To(Equal(uint64(123456789000000)))
		})

		It("unmarshals the JSON number to the expected uint", func() {
			u, err := models.UnmarshalUint(json.Number("123456789000000"))
			Expect(err).ToNot(HaveOccurred())
			Expect(u).To(Equal(uint64(123456789000000)))
		})

		It("unmarshals the integers to the expected uint", func() {
			u, err := models.UnmarshalUint(123456789)
			Expect(err).ToNot(HaveOccurred())
			Expect(u).To(Equal(uint64(123456789)))
		})

		It("errors if the data is not an integer type", func() {
			_, err := models.UnmarshalUint(1.2)
			Expect(err).To(MatchError(MatchRegexp(`.* is not a supported integer type`)))
		})
	})

})
