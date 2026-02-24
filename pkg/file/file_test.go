// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package file

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/check.v1"

	"github.com/stretchr/testify/require"

	"github.com/moov-io/metro2/pkg/lib"
	"github.com/moov-io/metro2/pkg/utils"
)

func TestFile__Crashers(t *testing.T) {
	paths := readCrasherInputFilePaths(t)
	for i := range paths {

		f, err := os.Open(paths[i])
		if err != nil {
			t.Fatal(err)
		}

		if testing.Verbose() {
			t.Logf("parsing %s", paths[i])
		}

		if _, err := NewFileFromReader(f); err == nil {
			t.Errorf("expected error with %s", paths[i])
		} else {
			t.Logf("error with %s\n  %#v", paths[i], err)
		}

		if testing.Verbose() {
			t.Logf("read %s without crashing", paths[i])
		}

		f.Close()
	}
}

func readCrasherInputFilePaths(t *testing.T) []string {
	t.Helper()

	basePath := filepath.Join("..", "..", "test", "testdata", "crashers")
	fds, err := os.ReadDir(basePath)
	if err != nil {
		t.Fatal(err)
	}

	var out []string
	for i := range fds {
		if strings.HasSuffix(fds[i].Name(), ".output") {
			continue
		}
		out = append(out, filepath.Join(basePath, fds[i].Name()))
	}
	return out
}

func (t *FileTest) TestJsonWithUnpackedVariableBlocked(c *check.C) {
	f, err := NewFile(utils.CharacterFileFormat)
	c.Assert(err, check.IsNil)
	err = json.Unmarshal(t.unpackedVariableBlockedJson, f)
	c.Assert(err, check.IsNil)

	raw, err := os.ReadFile(filepath.Join("..", "..", "test", "testdata", "unpacked_variable_file.dat"))
	c.Assert(err, check.IsNil)

	rawStr := strings.ReplaceAll(string(raw), "\r\n", "\n")
	fileString := f.String(true)
	compare := strings.Compare(fileString, rawStr)
	c.Assert(compare, check.Equals, 0)
	c.Assert(strings.Compare(f.ConcurrentString(true, 2), rawStr), check.Equals, 0)

	buf, err := json.Marshal(f)
	c.Assert(err, check.IsNil)
	var out bytes.Buffer
	err = json.Indent(&out, buf, "", "  ")
	c.Assert(err, check.IsNil)
	jsonStr := out.String()
	jsonStr = strings.ReplaceAll(jsonStr, "\n", "")
	c.Assert(jsonStr, check.Equals, string(t.unpackedVariableBlockedJson))
}

func (t *FileTest) TestJsonWithUnpackedFixedLength(c *check.C) {
	f, err := NewFile(utils.CharacterFileFormat)
	c.Assert(err, check.IsNil)
	err = json.Unmarshal(t.unpackedFixedLengthJson, f)
	c.Assert(err, check.IsNil)
	buf, err := json.Marshal(f)
	c.Assert(err, check.IsNil)
	var out bytes.Buffer
	err = json.Indent(&out, buf, "", "  ")
	c.Assert(err, check.IsNil)
	jsonStr := out.String()
	jsonStr = strings.ReplaceAll(jsonStr, "\n", "")
	c.Assert(jsonStr, check.Equals, string(t.unpackedFixedLengthJson))
}

func (t *FileTest) TestParseWithUnpackedFixedLength(c *check.C) {
	f, err := NewFile(utils.CharacterFileFormat)
	c.Assert(err, check.IsNil)
	err = f.Parse(t.unpackedFixedLengthRaw, false)
	c.Assert(err, check.IsNil)
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	trailer := f.GetDataRecords()[0]
	a, ok := trailer.(*lib.BaseSegment)
	c.Assert(ok, check.Equals, true)
	a.AccountStatus = lib.AccountStatusDF
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatusDA
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus05
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus11
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus13
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus61
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus63
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus64
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus65
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
}

func (t *FileTest) TestParseWithUnpackedFixedLength2(c *check.C) {
	f, err := NewFile(utils.CharacterFileFormat)
	c.Assert(err, check.IsNil)
	err = f.Parse(t.unpackedFixedLengthRaw, false)
	c.Assert(err, check.IsNil)
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	trailer := f.GetDataRecords()[0]
	a, ok := trailer.(*lib.BaseSegment)
	c.Assert(ok, check.Equals, true)
	a.AccountStatus = lib.AccountStatus71
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus78
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus80
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus82
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus83
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus84
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus88
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus89
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus93
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus94
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus95
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus96
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	a.AccountStatus = lib.AccountStatus97
	_, err = f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
}

func (t *FileTest) TestJsonWithPackedBlocked(c *check.C) {
	f, err := NewFile(utils.PackedFileFormat)
	c.Assert(err, check.IsNil)
	err = json.Unmarshal(t.packedJson, f)
	c.Assert(err, check.IsNil)
	buf, err := json.Marshal(f)
	c.Assert(err, check.IsNil)
	var out bytes.Buffer
	err = json.Indent(&out, buf, "", "  ")
	c.Assert(err, check.IsNil)
	jsonStr := out.String()
	jsonStr = strings.ReplaceAll(jsonStr, "\n", "")
	c.Assert(jsonStr, check.Equals, string(t.packedJson))
}

func (t *FileTest) TestParseWithPackedFileParse(c *check.C) {
	f, err := NewFile(utils.PackedFileFormat)
	c.Assert(err, check.IsNil)
	err = f.Parse(t.packedRaw, true)
	c.Assert(err, check.IsNil)
}

func (t *FileTest) TestFileSetBlock(c *check.C) {
	f, err := NewFile(utils.CharacterFileFormat)
	c.Assert(err, check.IsNil)
	jsonStr := `{
		"blockDescriptorWord": 430,
		"recordDescriptorWord": 426,
		"recordIdentifier": "HEADER",
		"transUnionProgramIdentifier": "5555555555",
		"activityDate": "2002-08-20T00:00:00Z",
		"dateCreated": "1999-05-10T00:00:00Z",
		"programDate": "1999-05-10T00:00:00Z",
		"programRevisionDate": "1999-05-10T00:00:00Z",
		"reporterName": "YOUR BUSINESS NAME HERE",
		"reporterAddress": "LINE ONE OF YOUR ADDRESS LINE TWO OF YOUR ADDRESS LINE THERE OF YOUR ADDRESS",
		"reporterTelephoneNumber": 1234567890
	  }`
	newSegment := lib.HeaderRecord{}
	err = json.Unmarshal([]byte(jsonStr), &newSegment)
	c.Assert(err, check.IsNil)
	orgHeader, err := f.GetRecord(utils.HeaderRecordName)
	c.Assert(err, check.IsNil)
	origin := orgHeader.BlockSize()
	err = f.SetRecord(&newSegment)
	c.Assert(err, check.IsNil)
	newHeader, err := f.GetRecord(utils.HeaderRecordName)
	c.Assert(err, check.IsNil)
	c.Assert(origin, check.Not(check.Equals), newHeader.BlockSize())
}

func (t *FileTest) TestFileDataRecord(c *check.C) {
	f, err := NewFile(utils.CharacterFileFormat)
	c.Assert(err, check.IsNil)
	segment := lib.NewBaseSegment()
	err = json.Unmarshal(t.baseSegmentJson, segment)
	c.Assert(err, check.IsNil)
	err = f.AddDataRecord(segment)
	c.Assert(err, check.IsNil)
	list := f.GetDataRecords()
	c.Assert(len(list), check.Equals, 1)
}

func (t *FileTest) TestGeneratorTrailer(c *check.C) {
	f, err := NewFile(utils.CharacterFileFormat)
	c.Assert(err, check.IsNil)
	err = json.Unmarshal(t.unpackedFixedLengthJson, f)
	c.Assert(err, check.IsNil)
	trailer, err := f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	err = f.SetRecord(trailer)
	c.Assert(err, check.IsNil)
	err = f.Validate()
	c.Assert(err, check.IsNil)
}

func (t *FileTest) TestGeneratorPackedTrailer(c *check.C) {
	f, err := NewFile(utils.PackedFileFormat)
	c.Assert(err, check.IsNil)
	err = json.Unmarshal(t.packedJson, f)
	c.Assert(err, check.IsNil)
	trailer, err := f.GeneratorTrailer()
	c.Assert(err, check.IsNil)
	err = f.SetRecord(trailer)
	c.Assert(err, check.IsNil)
	err = f.Validate()
	c.Assert(err, check.IsNil)
}

func (t *FileTest) TestGeneratorTrailerStatistics(c *check.C) {
	t.runTrailerTest(c, utils.CharacterFileFormat, &lib.BaseSegment{})
}

func (t *FileTest) TestGeneratorPackedTrailerStatistics(c *check.C) {
	t.runTrailerTest(c, utils.PackedFileFormat, &lib.PackedBaseSegment{})
}

// allAccountStatusCases returns every account status code and the trailer field that should be 1.
func allAccountStatusCases() []string {
	return []string{
		lib.AccountStatusDF,
		lib.AccountStatusDA,
		lib.AccountStatus05,
		lib.AccountStatus11,
		lib.AccountStatus13,
		lib.AccountStatus61,
		lib.AccountStatus62,
		lib.AccountStatus63,
		lib.AccountStatus64,
		lib.AccountStatus65,
		lib.AccountStatus71,
		lib.AccountStatus78,
		lib.AccountStatus80,
		lib.AccountStatus82,
		lib.AccountStatus83,
		lib.AccountStatus84,
		lib.AccountStatus88,
		lib.AccountStatus89,
		lib.AccountStatus93,
		lib.AccountStatus94,
		lib.AccountStatus95,
		lib.AccountStatus96,
		lib.AccountStatus97,
	}
}

// runTrailerTest is a helper function to test trailer generation for both character and packed file
// formats using the same logic. It tests for correct counting of account status codes and segment
// totals in the generated trailer.
func (t *FileTest) runTrailerTest(c *check.C, format string, baseTemplate interface{}) {
	f, err := os.Open(filepath.Join("..", "..", "test", "testdata", "trailer_statistics.json"))
	c.Assert(err, check.IsNil)
	defer f.Close()

	file, err := NewFile(format)
	c.Assert(err, check.IsNil)

	base := reflect.New(reflect.TypeOf(baseTemplate).Elem()).Interface()
	err = json.Unmarshal(utils.ReadFile(f), base) // <--- THIS WAS MISSING
	c.Assert(err, check.IsNil)
	baseVal := reflect.ValueOf(base).Elem()

	statusesRequiringPaymentRating := map[string]bool{
		lib.AccountStatus05: true, lib.AccountStatus13: true, lib.AccountStatus65: true,
		lib.AccountStatus88: true, lib.AccountStatus89: true, lib.AccountStatus94: true, lib.AccountStatus95: true,
	}

	// create a base segment for each account status, with corresponding payment rating
	accountStatuses := allAccountStatusCases()
	for _, status := range accountStatuses {
		baseCopyVal := reflect.New(baseVal.Type())
		baseCopyVal.Elem().Set(baseVal)

		v := baseCopyVal.Elem()

		v.FieldByName("AccountStatus").SetString(status)

		if statusesRequiringPaymentRating[status] {
			v.FieldByName("PaymentRating").SetString(lib.PaymentRatingCurrent)
		}

		record := baseCopyVal.Interface().(lib.Record)
		err = file.AddDataRecord(record)
		c.Assert(err, check.IsNil)
	}

	tr, err := file.GeneratorTrailer()
	c.Assert(err, check.IsNil)

	t.assertTrailerFields(c, tr, len(accountStatuses))
}

func (t *FileTest) assertTrailerFields(c *check.C, trailer any, numSegments int) {
	v := reflect.ValueOf(trailer).Elem()

	checkField := func(fieldName string, expected int) {
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			c.Fatalf("Field %s not found on trailer struct", fieldName)
		}
		actual := int(field.Int())
		c.Check(actual, check.Equals, expected,
			check.Commentf("Field %s: expected %d, got %d", fieldName, expected, actual))
	}

	accountStatusFields := []string{
		"TotalStatusCodeDF", "TotalStatusCodeDA", "TotalStatusCode05", "TotalStatusCode11",
		"TotalStatusCode13", "TotalStatusCode61", "TotalStatusCode62", "TotalStatusCode63",
		"TotalStatusCode64", "TotalStatusCode65", "TotalStatusCode71", "TotalStatusCode78",
		"TotalStatusCode80", "TotalStatusCode82", "TotalStatusCode83", "TotalStatusCode84",
		"TotalStatusCode88", "TotalStatusCode89", "TotalStatusCode93", "TotalStatusCode94",
		"TotalStatusCode95", "TotalStatusCode96", "TotalStatusCode97",
	}
	for _, f := range accountStatusFields {
		checkField(f, 1)
	}

	segmentTotalFields := []string{
		"TotalBaseRecords", "TotalConsumerSegmentsJ1", "TotalConsumerSegmentsJ2",
		"TotalOriginalCreditorSegments", "TotalPurchasedToSegments",
		"TotalMortgageInformationSegments", "TotalPaymentInformationSegments",
		"TotalChangeSegments", "TotalEmploymentSegments", "TotalSocialNumbersBaseSegments",
		"TotalSocialNumbersJ1Segments", "TotalSocialNumbersJ2Segments",
		"TotalDatesBirthBaseSegments", "TotalDatesBirthJ1Segments", "TotalDatesBirthJ2Segments",
	}
	for _, f := range segmentTotalFields {
		checkField(f, numSegments)
	}

	checkField("BlockCount", numSegments+2)

	// check combined segment fields
	checkField("TotalECOACodeZ", 3*numSegments)
	checkField("TotalSocialNumbersAllSegments", 3*numSegments)
	checkField("TotalTelephoneNumbersAllSegments", 3*numSegments)
}

func (t *FileTest) TestFileValidate(c *check.C) {
	f, err := NewFile(utils.PackedFileFormat)
	c.Assert(err, check.IsNil)
	err = json.Unmarshal(t.packedJson, f)
	c.Assert(err, check.IsNil)
	err = f.Validate()
	c.Assert(err, check.IsNil)
}

func (t *FileTest) TestGetRecord(c *check.C) {
	f, err := NewFile(utils.PackedFileFormat)
	c.Assert(err, check.IsNil)
	err = json.Unmarshal(t.packedJson, f)
	c.Assert(err, check.IsNil)
	_, err = f.GetRecord(lib.TrailerRecordName)
	c.Assert(err, check.IsNil)
	_, err = f.GetRecord(lib.BaseSegmentName)
	c.Assert(err, check.NotNil)
}

func (t *FileTest) TestCreateFile(c *check.C) {
	_, err := CreateFile(t.packedJson)
	c.Assert(err, check.IsNil)

	f, err := CreateFile(t.packedRaw)
	c.Assert(err, check.IsNil)

	raw, err := os.ReadFile(filepath.Join("..", "..", "test", "testdata", "packed_file.dat"))
	c.Assert(err, check.IsNil)

	c.Assert(strings.Compare(f.String(false), string(raw)), check.Equals, 0)
	c.Assert(strings.Compare(f.ConcurrentString(false, 2), string(raw)), check.Equals, 0)
}

func (t *FileTest) TestNewFileFromReader(c *check.C) {
	_, err := NewFileFromReader(t.packedJsonReader)
	c.Assert(err, check.IsNil)

	f, err := NewFileFromReader(t.packedRawReader)
	c.Assert(err, check.IsNil)

	raw, err := os.ReadFile(filepath.Join("..", "..", "test", "testdata", "packed_file.dat"))
	c.Assert(err, check.IsNil)

	c.Assert(strings.Compare(f.String(false), string(raw)), check.Equals, 0)
	c.Assert(strings.Compare(f.ConcurrentString(false, 2), string(raw)), check.Equals, 0)
}

func (t *FileTest) TestCreateFileFailed(c *check.C) {
	r1 := bytes.NewReader(t.packedRaw[8:])
	c.Assert(r1, check.NotNil)

	_, err := NewFileFromReader(r1)
	c.Assert(err, check.NotNil)

	data := `{
  "header": {
    "recordDescriptorWord": 480,
    "recordIdentifier": "error",
  }
}`
	r2 := bytes.NewReader([]byte(data))

	_, err = NewFileFromReader(r2)
	c.Assert(err, check.NotNil)
}

func (t *FileTest) TestWithUnknownFileType(c *check.C) {
	_, err := NewFile("unknown")
	c.Assert(err, check.NotNil)
}

func TestFile__Reader(t *testing.T) {
	t.Run("Read with unpacked fixed file", func(t *testing.T) {
		readUnpackedFixedFile(t)
	})

	t.Run("Read with packed file", func(t *testing.T) {
		readPackedFile(t)
	})
}

func BenchmarkFile(b *testing.B) {
	b.Run("Read with unpacked fixed file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			readUnpackedFixedFile(b)
		}
	})

	b.Run("Read with unpacked variable file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			readUnpackedVariableFile(b)
		}
	})

	b.Run("Read with packed file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			readPackedFile(b)
		}
	})
}

func readUnpackedFixedFile(tb testing.TB) {
	tb.Helper()

	fd, err := os.Open(filepath.Join("..", "..", "test", "testdata", "unpacked_fixed_file.dat"))
	if err != nil {
		tb.Fatalf("Can not open local file: %s: \n", err)
	}
	defer fd.Close()

	f, err := NewReader(fd).Read()
	require.NoError(tb, err)

	// ensure we have a validated file structure
	err = f.Validate()
	require.NoError(tb, err)
}

func readUnpackedVariableFile(tb testing.TB) {
	tb.Helper()

	fd, err := os.Open(filepath.Join("..", "..", "test", "testdata", "unpacked_variable_file.dat"))
	if err != nil {
		tb.Fatalf("Can not open local file: %s: \n", err)
	}
	defer fd.Close()

	f, err := NewReader(fd).Read()
	require.NoError(tb, err)

	// ensure we have a validated file structure
	err = f.Validate()
	require.NoError(tb, err)
}

func readPackedFile(tb testing.TB) {
	tb.Helper()

	fd, err := os.Open(filepath.Join("..", "..", "test", "testdata", "packed_file.dat"))
	if err != nil {
		tb.Fatalf("Can not open local file: %s: \n", err)
	}
	defer fd.Close()

	f, err := NewReader(fd).Read()
	require.NoError(tb, err)

	// ensure we have a validated file structure
	err = f.Validate()
	require.NoError(tb, err)
}
