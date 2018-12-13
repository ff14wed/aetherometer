// +build ignore

package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/ff14wed/sibyl/backend/datasheet"
)

type DataEntry map[string]interface{}

type Datasheet []DataEntry

type Schema struct {
	Names []string
	Types []string
}

func parseCSV(dataBytes io.Reader) (Datasheet, error) {
	csvReader := csv.NewReader(dataBytes)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 4 {
		return nil, errors.New("no records in data sheet")
	}

	schema := Schema{
		Names: make([]string, len(records[0])),
		Types: make([]string, len(records[0])),
	}

	for i, entry := range records[0] {
		schema.Names[i] = entry
	}
	for i, entry := range records[1] {
		if entry != "" && entry != "#" {
			schema.Names[i] = entry
		}
	}
	schema.Names[0] = "key"
	for i, entry := range records[2] {
		schema.Types[i] = entry
	}

	d := make(Datasheet, len(records[3:]))

	for i, record := range records[3:] {
		d[i] = make(DataEntry)
		for j, entry := range record {
			d[i][schema.Names[j]] = convertType(schema.Types[j], entry)
		}
	}
	return d, nil
}

func convertType(typeName string, entry string) interface{} {
	switch strings.ToLower(typeName) {
	case "byte", "sbyte", "uint8", "uint16", "uint32", "uint64",
		"int8", "int16", "int32", "int64":
		v, err := strconv.Atoi(entry)
		if err != nil {
			return entry
		}
		return v
	case "single":
		v, err := strconv.ParseFloat(entry, 64)
		if err != nil {
			return entry
		}
		return v
	}
	if strings.Contains(typeName, "bit&") {
		if (entry == "True") || (entry == "true") {
			return true
		}
		if (entry == "False") || (entry == "false") {
			return false
		}
	}
	return entry
}

func main() {
	csvFilenames := os.Args[1:]
	if len(csvFilenames) == 0 {
		log.Fatalf("%s [List of filenames to convert]: Must provide csv arguments\n", os.Args[0])
	}

	r := datasheet.FileReader{}
	var ds Datasheet

	for _, csvFilename := range csvFilenames {
		r.ReadFile(csvFilename, func(f io.Reader) error {
			var err error
			ds, err = parseCSV(f)
			if err != nil {
				return err
			}
			jsonBytes, err := json.MarshalIndent(ds, "", "  ")
			if err != nil {
				return err
			}
			ext := path.Ext(csvFilename)
			newFilename := csvFilename[0:len(csvFilename)-len(ext)] + ".json"
			return ioutil.WriteFile(newFilename, jsonBytes, os.ModePerm)
		})
	}
	if r.Error() != nil {
		log.Fatalf("Error processing files: %s\n", r.Error())
	}
	log.Printf("Successfully converted %d files\n", len(csvFilenames))
}
