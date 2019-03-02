package datasheet

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
)

// DataEntry is an internal structure for holding a single CSV record
type DataEntry map[string]string

// Datasheet is an internal structure for holding a set of CSV records
type Datasheet []DataEntry

// ParseRawCSV returns a raw Datasheet that represents the CSV data after
// parsing the table headers.
func ParseRawCSV(dataReader io.Reader) (Datasheet, error) {
	csvReader := csv.NewReader(dataReader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 4 {
		return nil, errors.New("no records in data sheet")
	}

	headers := make([]string, len(records[0]))

	copy(headers, records[0])

	for i, entry := range records[1] {
		if entry != "" && entry != "#" {
			headers[i] = entry
		}
	}
	headers[0] = "key"

	d := make(Datasheet, len(records[3:]))

	for i, record := range records[3:] {
		d[i] = make(DataEntry)
		for j, entry := range record {
			d[i][headers[j]] = entry
		}
	}
	return d, nil
}

// UnmarshalReader is a convenience wrapper around reading bytes from an
// io.Reader and storing the parsed data into the slice of structs pointed to
// by v.
func UnmarshalReader(dataReader io.Reader, v interface{}) error {
	dataBytes, err := ioutil.ReadAll(dataReader)
	if err != nil {
		return fmt.Errorf("error reading data: %s", err)
	}
	return Unmarshal(dataBytes, v)
}

// Unmarshal parses the CSV-encoded datasheet and stores the
// result in the slice of structs pointed to by v. It works very similarly to to
// json.Unmarshal.
//
// If v is nil or not a pointer to a slice of structs, Unmarshal returns an
// InvalidUnmarshalError.
// If the data is invalid CSV data, it returns a error from ParseRawCSV.
// If Unmarshal is unable to assign to a field of v the data, it returns
// an UnmarshalTypeError.
func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	rSlice := rv.Elem()
	if rSlice.Kind() != reflect.Slice {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	structType := rSlice.Type().Elem()
	if structType.Kind() != reflect.Struct {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	ds, err := ParseRawCSV(bytes.NewReader(data))
	if err != nil {
		return err
	}

	for i, source := range ds {
		v, err := convertStruct(structType, source, i)
		if err != nil {
			return err
		}
		rSlice.Set(reflect.Append(rSlice, v))
	}

	return nil
}

// convertStruct converts the data entry to the target struct type and
// returns the reflect.Value storing this converted data
func convertStruct(target reflect.Type, source DataEntry, row int) (reflect.Value, error) {
	numFields := target.NumField()
	targetVal := reflect.New(target).Elem()

	keyToFieldMapping := make(map[string]int)
	structTagFields := make(map[string]int)

	for i := 0; i < numFields; i++ {
		key := target.Field(i).Tag.Get("datasheet")
		if key != "" {
			if _, exists := structTagFields[key]; exists {
				return reflect.Value{}, &RepeatStructTagError{Field: target.Field(i).Name, Tag: key}
			}
			structTagFields[key] = i
			keyToFieldMapping[key] = i
			continue
		}
		key = target.Field(i).Name
		// Only fallback to mapping the key to a field name option if the key
		// isn't specified by a struct tag
		if _, exists := keyToFieldMapping[key]; !exists {
			keyToFieldMapping[key] = i
		}
	}
	for k, f := range keyToFieldMapping {
		sv, found := source[k]
		if !found {
			continue
		}
		fieldError := func() error {
			return &UnmarshalTypeError{
				Value:  sv,
				Type:   target.Field(f).Type,
				Row:    row,
				Struct: target.Name(),
				Field:  target.Field(f).Name,
			}
		}

		switch target.Field(f).Type.Kind() {
		case reflect.Bool:
			if (sv == "True") || (sv == "true") {
				targetVal.Field(f).SetBool(true)
			}
		case reflect.String:
			targetVal.Field(f).SetString(sv)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.Atoi(sv)
			if err != nil {
				return reflect.Value{}, fieldError()
			}
			targetVal.Field(f).SetUint(uint64(v))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.Atoi(sv)
			if err != nil {
				return reflect.Value{}, fieldError()
			}
			targetVal.Field(f).SetInt(int64(v))
		case reflect.Float32:
			v, err := strconv.ParseFloat(sv, 32)
			if err != nil {
				return reflect.Value{}, fieldError()
			}
			targetVal.Field(f).SetFloat(v)
		}
	}
	return targetVal, nil
}

// UnmarshalTypeError describes an error returned when the
// data could not be assigned to a specific Go type
type UnmarshalTypeError struct {
	Value  string
	Type   reflect.Type
	Row    int
	Struct string
	Field  string
}

func (e *UnmarshalTypeError) Error() string {
	return "datasheet: cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
}

// InvalidUnmarshalError describes an error returned when
// an invalid argument is passed to Unmarshal
type InvalidUnmarshalError struct {
	Type reflect.Type
}

// Error returns a formatted error message for the InvalidUnmarshalError
func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "datasheet: Unmarshal(nil)"
	}

	t := e.Type
	if t.Kind() != reflect.Ptr {
		return "datasheet: Unmarshal(non-pointer " + t.String() + ")"
	} else if t.Elem().Kind() != reflect.Slice || t.Elem().Elem().Kind() != reflect.Struct {
		return "datasheet: Unmarshal(not a pointer to a slice of structs: " + e.Type.String() + ")"
	}

	return "datasheet: Unmarshal(nil " + t.String() + ")"
}

// RepeatStructTagError describes an error returned when a repeat struct tag
// is encountered
type RepeatStructTagError struct {
	Field string
	Tag   string
}

// Error returns a formatted error message for the RepeatStructTagError
func (e *RepeatStructTagError) Error() string {
	return "datasheet: struct field " + e.Field + " repeats datasheet struct tag " + e.Tag
}
