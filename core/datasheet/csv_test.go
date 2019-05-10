package datasheet_test

import (
	"bytes"

	"github.com/ff14wed/aetherometer/core/datasheet"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testCSVData = `
key,0,1,2,3,4,5,6,7
#,Singular,,Plural,,,,,
int32,str,sbyte,str,sbyte,sbyte,sbyte,single,bit&01
0,"",0,"",0,0,1,0,False
1,"",0,"",0,0,0,0,False
2,"ruins runner",0,"ruins runners",0,0,1,0.1,True
3,"antelope doe",1,"antelope does",0,1,1,0.2,True
`

var expectedDatasheet = datasheet.Datasheet{
	datasheet.DataEntry{
		"key": "0", "Singular": "", "1": "0", "Plural": "", "3": "0",
		"4": "0", "5": "1", "6": "0", "7": "False",
	},
	datasheet.DataEntry{
		"key": "1", "Singular": "", "1": "0", "Plural": "", "3": "0",
		"4": "0", "5": "0", "6": "0", "7": "False",
	},
	datasheet.DataEntry{
		"key": "2", "Singular": "ruins runner", "1": "0",
		"Plural": "ruins runners", "3": "0", "4": "0", "5": "1", "6": "0.1", "7": "True",
	},
	datasheet.DataEntry{
		"key": "3", "Singular": "antelope doe", "1": "1",
		"Plural": "antelope does", "3": "0", "4": "1", "5": "1", "6": "0.2", "7": "True",
	},
}

type testDataStruct struct {
	ID           uint32 `datasheet:"key"`
	Singular     string
	Extra        byte `datasheet:"1"`
	Plural       string
	FakeFloat    float32 `datasheet:"6"`
	FakeBool     bool    `datasheet:"7"`
	Nonexistent  int
	Nonexistent2 float32 `datasheet:"nonexistent2"`
}

var expectedParsedData = []testDataStruct{
	testDataStruct{
		ID: 0, Singular: "", Extra: 0, Plural: "", FakeFloat: 0, FakeBool: false,
	},
	testDataStruct{
		ID: 1, Singular: "", Extra: 0, Plural: "", FakeFloat: 0, FakeBool: false,
	},
	testDataStruct{
		ID: 2, Singular: "ruins runner", Extra: 0, Plural: "ruins runners", FakeFloat: 0.1, FakeBool: true,
	},
	testDataStruct{
		ID: 3, Singular: "antelope doe", Extra: 1, Plural: "antelope does", FakeFloat: 0.2, FakeBool: true,
	},
}

var _ = Describe("CSV", func() {
	Describe("ParseRawCSV", func() {
		It("returns an array of data records with the names labeled according to the sheet headers", func() {
			ds, err := datasheet.ParseRawCSV(bytes.NewBufferString(testCSVData))
			Expect(err).ToNot(HaveOccurred())
			Expect(ds).To(Equal(expectedDatasheet))
		})

		Context("when the csv data is invalid", func() {
			const invalidData = "key,0,1\n#,Singular"

			It("returns an error", func() {
				_, err := datasheet.ParseRawCSV(bytes.NewBufferString(invalidData))
				Expect(err).To(MatchError(ContainSubstring("wrong number of fields")))
			})
		})

		Context("when there aren't any records", func() {
			const invalidData = "key,0,1\n#,Singular,\nint32,str,sbyte"

			It("returns an error", func() {
				_, err := datasheet.ParseRawCSV(bytes.NewBufferString(invalidData))
				Expect(err).To(MatchError(ContainSubstring("no records in data sheet")))
			})
		})
	})

	Describe("Unmarshal", func() {
		It("successfully unmarshals the data into the provided slice of structs", func() {
			var v []testDataStruct
			Expect(datasheet.Unmarshal([]byte(testCSVData), &v)).To(Succeed())
			Expect(v).To(Equal(expectedParsedData))
		})

		Context("when the data does not match the struct type", func() {
			const testCSV = `key,0,1,2,3,4,5,6,7` + "\n" +
				`#,Singular,,Nonexistent,,,,,` + "\n" +
				`int32,str,sbyte,str,sbyte,sbyte,sbyte,single,bit&01` + "\n" +
				`0,"foo",0,"foos",0,0,1,0,False`

			It("returns an UnmarshalTypeError", func() {
				var v []testDataStruct
				err := datasheet.Unmarshal([]byte(testCSV), &v)
				Expect(err).To(MatchError("datasheet: cannot unmarshal foos into Go struct field testDataStruct.Nonexistent of type int"))
			})
		})

		Context("when both a struct tag and a separate field name reference a data header", func() {
			type ambiguousStruct struct {
				Field     string
				RealField string `datasheet:"Field"`
			}

			const testCSV = `key,0` + "\n" +
				`#,Field` + "\n" +
				`int32,str` + "\n" +
				`0,"foo"`

			It("the struct tag wins", func() {
				var v []ambiguousStruct
				err := datasheet.Unmarshal([]byte(testCSV), &v)
				Expect(err).ToNot(HaveOccurred())
				Expect(v).To(ConsistOf(ambiguousStruct{
					RealField: "foo",
				}))
			})
		})

		Context("when a struct tag repeats", func() {
			type ambiguousStruct struct {
				Field     string
				FakeField string `datasheet:"Field"`
				RealField string `datasheet:"Field"`
			}

			const testCSV = `key,0` + "\n" +
				`#,Field` + "\n" +
				`int32,str` + "\n" +
				`0,"foo"`

			It("returns a RepeatStructTagError", func() {
				var v []ambiguousStruct
				err := datasheet.Unmarshal([]byte(testCSV), &v)
				Expect(err).To(MatchError("datasheet: struct field RealField repeats datasheet struct tag Field"))
			})

		})

		It("successfully unmarshals the data into the provided slice of structs", func() {
			var v []testDataStruct
			Expect(datasheet.Unmarshal([]byte(testCSVData), &v)).To(Succeed())
			Expect(v).To(Equal(expectedParsedData))
		})

		Describe("when the output value is an invalid type", func() {
			It("returns an error when the value is nil", func() {
				err := datasheet.Unmarshal([]byte(testCSVData), nil)
				Expect(err).To(MatchError("datasheet: Unmarshal(nil)"))
			})

			It("returns an error when the value is a non-pointer", func() {
				var v []testDataStruct
				err := datasheet.Unmarshal([]byte(testCSVData), v)
				Expect(err).To(MatchError("datasheet: Unmarshal(non-pointer []datasheet_test.testDataStruct)"))
			})

			It("returns an error when the value is a nil pointer", func() {
				var v *[]testDataStruct
				err := datasheet.Unmarshal([]byte(testCSVData), v)
				Expect(err).To(MatchError("datasheet: Unmarshal(nil *[]datasheet_test.testDataStruct)"))
			})

			It("returns an error when the value is a pointer to a non-struct", func() {
				v := 7
				err := datasheet.Unmarshal([]byte(testCSVData), &v)
				Expect(err).To(MatchError("datasheet: Unmarshal(not a pointer to a slice of structs: *int)"))
			})

			It("returns an error when the value is a pointer to a slice of non-struct", func() {
				var v []int
				err := datasheet.Unmarshal([]byte(testCSVData), &v)
				Expect(err).To(MatchError("datasheet: Unmarshal(not a pointer to a slice of structs: *[]int)"))
			})
		})

		Context("when the csv data is invalid", func() {
			const invalidData = "key,0,1\n#,Singular"

			It("returns an error", func() {
				var v []testDataStruct
				err := datasheet.Unmarshal([]byte(invalidData), &v)
				Expect(err).To(MatchError(ContainSubstring("wrong number of fields")))
			})
		})
	})
})
