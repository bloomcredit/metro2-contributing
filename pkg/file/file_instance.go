// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"sync"
	"unicode"

	"github.com/moov-io/base/log"
	"github.com/moov-io/metro2/pkg/lib"
	"github.com/moov-io/metro2/pkg/utils"
)

var _ File = (*fileInstance)(nil)

// File contains the structures of a parsed metro 2 file.
type fileInstance struct {
	Header  lib.Record   `json:"header"`
	Bases   []lib.Record `json:"data"`
	Trailer lib.Record   `json:"trailer"`

	format string
	logger log.Logger
}

// SetRecord can set block record like as header, trailer
func (f *fileInstance) SetRecord(r lib.Record) error {
	err := r.Validate()
	if err != nil {
		return err
	}

	if (f.format == utils.PackedFileFormat && r.Name() == lib.PackedHeaderRecordName) ||
		(f.format == utils.CharacterFileFormat && r.Name() == lib.HeaderRecordName) {
		f.Header = r
	} else if (f.format == utils.PackedFileFormat && r.Name() == lib.PackedTrailerRecordName) ||
		(f.format == utils.CharacterFileFormat && r.Name() == lib.TrailerRecordName) {
		f.Trailer = r
	} else {
		return utils.NewErrInvalidRecord(r.Name())
	}

	return nil
}

// AddDataRecord can append data record
func (f *fileInstance) AddDataRecord(r lib.Record) error {
	err := r.Validate()
	if err != nil {
		return err
	}

	if f.format == utils.PackedFileFormat && r.Name() == lib.PackedBaseSegmentName {
		f.Bases = append(f.Bases, r)
	} else if f.format == utils.CharacterFileFormat && r.Name() == lib.BaseSegmentName {
		f.Bases = append(f.Bases, r)
	} else {
		return utils.NewErrInvalidRecord(r.Name())
	}

	return nil
}

// GetRecord returns single record like as header, trailer.
func (f *fileInstance) GetRecord(name string) (lib.Record, error) {
	switch name {
	case utils.HeaderRecordName:
		return f.Header, nil
	case utils.TrailerRecordName:
		return f.Trailer, nil
	default:
		return nil, utils.NewErrInvalidRecord(name)
	}
}

// GetDataRecords returns data records
func (f *fileInstance) GetDataRecords() []lib.Record {
	return f.Bases
}

// GeneratorTrailer returns trailer segment that created automatically
func (f *fileInstance) GeneratorTrailer() (lib.Record, error) {
	var trailer lib.Record
	var err error

	if f.format == utils.PackedFileFormat {
		trailer, err = f.generatorPackedTrailer()
	} else {
		trailer, err = f.generatorTrailer()
	}
	if err != nil {
		return nil, err
	}

	return trailer, nil
}

// Validate performs some checks on the file and returns an error if not Validated
func (f *fileInstance) Validate() error {
	var err error

	if f.Header != nil {
		err = f.Header.Validate()
		if err != nil {
			return err
		}
	}

	if f.Trailer != nil {
		err = f.Trailer.Validate()
		if err != nil {
			return err
		}
	}

	for _, baseSegment := range f.Bases {
		err = baseSegment.Validate()
		if err != nil {
			return err
		}
	}

	var fromFields reflect.Value
	if f.format == utils.PackedFileFormat {
		trailer, err := f.generatorPackedTrailer()
		if err != nil {
			return err
		}
		fromFields = reflect.ValueOf(trailer).Elem()
	} else {
		trailer, err := f.generatorTrailer()
		if err != nil {
			return err
		}
		fromFields = reflect.ValueOf(trailer).Elem()
	}

	toFields := reflect.ValueOf(f.Trailer).Elem()
	for i := 0; i < fromFields.NumField(); i++ {
		fieldName := fromFields.Type().Field(i).Name

		// skip local variable
		if !unicode.IsUpper([]rune(fieldName)[0]) {
			continue
		}
		switch fieldName {
		case "BlockDescriptorWord", "RecordDescriptorWord", "RecordIdentifier":
			continue
		}

		fromField := fromFields.FieldByName(fieldName)
		toField := toFields.FieldByName(fieldName)
		if !fromField.IsValid() || !toField.IsValid() {
			return utils.NewErrInvalidValueOfField(fieldName, "trailer record")
		}
		if fromField.Interface() != toField.Convert(fromField.Type()).Interface() {
			return utils.NewErrInvalidValueOfField(fieldName, "trailer record")
		}
	}
	return nil
}

// Parse attempts to initialize a *File object assuming the input is valid raw data.
func (f *fileInstance) Parse(record []byte, isVariableLength bool) error {
	// remove new lines
	record = slices.DeleteFunc(record, func(b byte) bool {
		return b == '\r' || b == '\n'
	})

	f.Bases = []lib.Record{}
	offset := 0

	// Header Record
	head, err := f.Header.Parse(record, isVariableLength)
	if err != nil {
		return err
	}
	offset += head

	// Data Record
	for err == nil {
		var base lib.Record
		if f.format == utils.PackedFileFormat {
			base = lib.NewPackedBaseSegment()
		} else {
			base = lib.NewBaseSegment()
		}

		if offset <= 0 || len(record) <= offset {
			return utils.NewErrSegmentLength("base record")
		}

		read, err := base.Parse(record[offset:], isVariableLength)
		if err != nil {
			break
		}
		f.Bases = append(f.Bases, base)
		offset += read
	}

	// Trailer Record
	if offset <= 0 || len(record) <= offset {
		return utils.NewErrSegmentLength("trailer record")
	}
	tread, err := f.Trailer.Parse(record[offset:], isVariableLength)
	if err != nil {
		return err
	}
	offset += tread

	if offset != len(record) {
		return utils.NewErrFailedParsing()
	}

	return nil
}

var (
	// defaultStringConcurrency is the default number of goroutines to use for File.String() calls,
	// which is the number of CPUs detectable by Go.
	defaultStringConcurrency = runtime.NumCPU()
)

// String writes the File struct to raw string.
func (f *fileInstance) String(isNewLine bool) string {
	return f.ConcurrentString(isNewLine, defaultStringConcurrency)
}

// ConcurrentString writes the File struct to a string by concurrently generating rows.
func (f *fileInstance) ConcurrentString(isNewLine bool, goroutines int) string {
	if goroutines < 1 {
		goroutines = 1
	}

	newLine := ""
	if isNewLine {
		newLine = "\n"
	}

	// Header Block
	header := f.Header.String() + newLine

	// Data Block
	pageSize := int(math.Ceil(float64(len(f.Bases)) / float64(goroutines)))
	basePages := [][]lib.Record{}
	for i := 0; i < len(f.Bases); i += pageSize {
		end := i + pageSize
		if end > len(f.Bases) {
			end = len(f.Bases)
		}
		basePages = append(basePages, f.Bases[i:end])
	}

	// Determine record length based on file format for better pre-allocation
	recordLength := lib.UnpackedRecordLength
	if f.format == utils.PackedFileFormat {
		recordLength = lib.PackedRecordLength
	}
	recordLength += len(newLine)

	dataPages := make([]string, len(basePages))
	var wg sync.WaitGroup
	for i, page := range basePages {
		wg.Add(1)
		go func(idx int, page []lib.Record) {
			defer wg.Done()
			var data strings.Builder
			data.Grow(len(page) * recordLength)
			for _, base := range page {
				data.WriteString(base.String())
				data.WriteString(newLine)
			}
			dataPages[idx] = data.String()
		}(i, page)
	}
	wg.Wait()

	// Trailer Block
	trailer := f.Trailer.String()

	// Combine Blocks
	var buf strings.Builder
	dataLength := 0
	for _, page := range dataPages {
		dataLength += len(page)
	}
	buf.Grow(len(header) + dataLength + len(trailer))

	buf.WriteString(header)
	for _, page := range dataPages {
		buf.WriteString(page)
	}
	buf.WriteString(trailer)

	return buf.String()
}

// Bytes return raw byte array
func (f *fileInstance) Bytes() []byte {
	return []byte(f.String(false))
}

// UnmarshalJSON parses a JSON blob
func (f *fileInstance) UnmarshalJSON(data []byte) error {
	dummy := make(map[string]interface{})
	err := json.Unmarshal(data, &dummy)
	if err != nil {
		return err
	}

	for name, record := range dummy {
		buf, err := json.Marshal(record)
		if err != nil {
			return err
		}

		switch name {
		case utils.HeaderRecordName:
			err = json.Unmarshal(buf, f.Header)
			if err != nil {
				f.logger.Error().LogErrorf(err.Error())
				return errors.New("Unable to parse input json file")
			}
		case utils.TrailerRecordName:
			err = json.Unmarshal(buf, f.Trailer)
			if err != nil {
				f.logger.Error().LogErrorf(err.Error())
				return errors.New("Unable to parse input json file")
			}
		case utils.DataRecordName:
			var list []interface{}
			err = json.Unmarshal(buf, &list)
			if err != nil {
				f.logger.Error().LogErrorf(err.Error())
				return errors.New("Unable to parse input json file")
			}
			for _, subSegment := range list {
				subBuf, err := json.Marshal(subSegment)
				if err != nil {
					return err
				}
				if f.format == utils.CharacterFileFormat {
					base := lib.NewBaseSegment()
					err = json.Unmarshal(subBuf, base)
					if err != nil {
						f.logger.Error().LogErrorf(err.Error())
						return errors.New("Unable to parse input json file")
					}
					f.Bases = append(f.Bases, base)
				} else {
					base := lib.NewPackedBaseSegment()
					err = json.Unmarshal(subBuf, base)
					if err != nil {
						f.logger.Error().LogErrorf(err.Error())
						return errors.New("Unable to parse input json file")
					}
					f.Bases = append(f.Bases, base)
				}
			}
		}
	}

	return nil
}

// GetType returns type of the file
func (f *fileInstance) GetType() string {
	return f.format
}

// GetType returns type of the file
func (f *fileInstance) SetType(newType string) error {
	if newType != utils.CharacterFileFormat && newType != utils.PackedFileFormat {
		return errors.New("invalid type")
	}

	if f.format == newType {
		return nil
	}

	if newType == utils.CharacterFileFormat {
		// convert header
		if f.Header != nil {
			rec, ok := f.Header.(*lib.PackedHeaderRecord)
			if !ok || rec == nil {
				return fmt.Errorf("unexpected PackedHeaderRecord: %T", f.Header)
			}
			newHeader := lib.HeaderRecord(*rec)
			f.Header = &newHeader
		}

		// convert bases
		var bases []lib.Record
		for _, _base := range f.Bases {
			rec, ok := _base.(*lib.PackedBaseSegment)
			if !ok || rec == nil {
				return fmt.Errorf("unexpected PackedBaseSegment: %T", _base)
			}
			newBase := lib.BaseSegment(*rec)
			bases = append(bases, &newBase)
		}
		f.Bases = bases

		// convert trailer
		if f.Trailer != nil {
			rec, ok := f.Trailer.(*lib.PackedTrailerRecord)
			if !ok || rec == nil {
				return fmt.Errorf("unexpected PackedTrailerRecord: %T", f.Trailer)
			}
			newTrailer := lib.TrailerRecord(*rec)
			f.Trailer = &newTrailer
		}
	} else if newType == utils.PackedFileFormat {
		// convert header
		if f.Header != nil {
			rec, ok := f.Header.(*lib.HeaderRecord)
			if !ok || rec == nil {
				return fmt.Errorf("unexpected HeaderRecord: %T", f.Header)
			}
			newHeader := lib.PackedHeaderRecord(*rec)
			f.Header = &newHeader
		}

		// convert bases
		var bases []lib.Record
		for _, _base := range f.Bases {
			rec, ok := _base.(*lib.BaseSegment)
			if !ok || rec == nil {
				return fmt.Errorf("unexpected BaseSegment: %T", _base)
			}
			newBase := lib.PackedBaseSegment(*rec)
			bases = append(bases, &newBase)
		}
		f.Bases = bases

		// convert trailer
		if f.Trailer != nil {
			rec, ok := f.Trailer.(*lib.TrailerRecord)
			if !ok || rec == nil {
				return fmt.Errorf("unexpected TrailerRecord: %T", f.Trailer)
			}
			newTrailer := lib.PackedTrailerRecord(*rec)
			f.Trailer = &newTrailer
		}
	}

	f.format = newType
	return nil
}

func (f *fileInstance) generatorTrailer() (*lib.TrailerRecord, error) {
	trailer := &lib.TrailerRecord{
		RecordDescriptorWord: lib.UnpackedRecordLength,
		RecordIdentifier:     lib.TrailerIdentifier,
		BlockCount:           2,
	}

	for _, base := range f.Bases {
		baseSegment, ok := base.(*lib.BaseSegment)
		if !ok && baseSegment.Validate() != nil {
			return nil, utils.NewErrInvalidSegment(baseSegment.Name())
		}
		trailer.TallyDataRecord(baseSegment)
	}

	return trailer, nil
}

func (f *fileInstance) generatorPackedTrailer() (*lib.PackedTrailerRecord, error) {
	trailer := &lib.PackedTrailerRecord{
		RecordDescriptorWord: lib.PackedRecordLength,
		RecordIdentifier:     lib.TrailerIdentifier,
		BlockCount:           2,
	}

	for _, base := range f.Bases {
		base, ok := base.(*lib.PackedBaseSegment)
		if !ok && base.Validate() != nil {
			return nil, utils.NewErrInvalidSegment(base.Name())
		}
		trailer.TallyDataRecord(base)
	}

	return trailer, nil
}
