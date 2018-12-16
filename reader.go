package main

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"bufio"

	"io"

	"encoding/json"
	"io/ioutil"

	"fmt"
	"strconv"

	"github.com/tealeg/xlsx"
	"github.com/urfave/cli"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Stringer interface {
	String(string) (string, error)
}

type NoopDecoder struct{}

func (d *NoopDecoder) String(s string) (string, error) {
	return s, nil
}

type Reader interface {
	Read() ([]string, error)
	Close() error
}

type CsvReader struct {
	cr       *csv.Reader
	f        *os.File
	counter  int
	startRow int
}

func (r *CsvReader) Read() ([]string, error) {
	if r.startRow > r.counter {
		r.counter++
		return nil, nil
	}
	r.counter++
	return r.cr.Read()
}

func (r *CsvReader) Close() error {
	return r.f.Close()
}

func newCsvReader(filename string, encoding string, mode string, start int) (*CsvReader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	var r *csv.Reader
	switch strings.ToUpper(encoding) {
	case "UTF8", "UTF-8":
		r = csv.NewReader(f)
	case "SHIFT-JIS", "SJIS":
		r = csv.NewReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
	case "EUC-JP", "EUCJP":
		r = csv.NewReader(transform.NewReader(f, japanese.EUCJP.NewDecoder()))
	}
	if mode == "tsv" {
		r.Comma = '\t'
	}
	r.LazyQuotes = true

	return &CsvReader{cr: r, f: f, startRow: start}, nil
}

type ExcelReader struct {
	xf       *xlsx.File
	xs       *xlsx.Sheet
	counter  int
	startRow int
	maxRow   int
}

func (r *ExcelReader) Read() ([]string, error) {
	if r.maxRow <= r.counter {
		return nil, io.EOF
	}
	if r.startRow > r.counter {
		r.counter++
		return nil, nil
	}
	row := r.xs.Rows[r.counter]
	values := make([]string, len(row.Cells))
	for i, cell := range row.Cells {
		values[i] = cell.Value
	}
	r.counter++
	return values, nil
}

func (r *ExcelReader) Close() error { return nil }

func newExcelReader(f string, sheet string, start int) (*ExcelReader, error) {
	xf, err := xlsx.OpenFile(f)
	if err != nil {
		return nil, err
	}
	for _, s := range xf.Sheets {
		if s.Name == sheet {

			return &ExcelReader{
				counter:  0,
				xf:       xf,
				xs:       s,
				startRow: start,
				maxRow:   s.MaxRow,
			}, nil
		}
	}
	return nil, errors.New("Sheet does not exists")
}

type FixWidthFileReader struct {
	f           *os.File
	s           *bufio.Scanner
	e           string
	byteNumbers []int
}

func (r *FixWidthFileReader) Read() ([]string, error) {
	if r.s.Scan() {
		var s Stringer
		switch strings.ToUpper(r.e) {
		case "SHIFT-JIS", "SJIS", "SHIFT_JIS":
			s = japanese.ShiftJIS.NewDecoder()
		default:
			s = &NoopDecoder{}
		}
		t := r.s.Text()
		start := 0
		values := make([]string, len(r.byteNumbers))
		for i, n := range r.byteNumbers {
			if len(t) < start+n {
				return nil, fmt.Errorf("dat is not valid length at line %d", i)
			}
			value, err := s.String(t[start : start+n])
			if err != nil {
				return nil, err
			}
			values[i] = strings.TrimSpace(value)
			start += n
		}
		return values, nil
	}
	return nil, nil
}

func (r *FixWidthFileReader) Close() error {
	return r.f.Close()
}

func newFixWidthFileReader(f string, e string, byteNumbers []int) (*FixWidthFileReader, error) {
	fp, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(fp)
	return &FixWidthFileReader{f: fp, s: s, e: e, byteNumbers: byteNumbers}, nil
}

type JsonReader struct {
	records  []map[string]interface{}
	f        *os.File
	counter  int
	startRow int
}

func (r *JsonReader) Read() ([]string, error) {
	if len(r.records) <= r.counter {
		return nil, io.EOF
	}
	if r.startRow > r.counter {
		r.counter++
		return nil, nil
	}
	record := r.records[r.counter]
	values := make([]string, len(record))
	r.counter++
	return values, nil
}

func (r *JsonReader) Close() error {
	return r.f.Close()
}

func newJsonReader(filename string, startRow int) (*JsonReader, error) {
	records := []map[string]interface{}{}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(b, records)
	return &JsonReader{records: records, startRow: startRow}, nil
}

type JsonlReader struct {
	records  []map[string]interface{}
	f        *os.File
	counter  int
	startRow int
}

func (r *JsonlReader) Read() ([]string, error) {
	if len(r.records) <= r.counter {
		return nil, io.EOF
	}
	if r.startRow > r.counter {
		r.counter++
		return nil, nil
	}
	record := r.records[r.counter]
	values := make([]string, len(record))
	r.counter++
	return values, nil
}

func (r *JsonlReader) Close() error {
	return r.f.Close()
}

func newJsonlReader(filename string, startRow int) (*JsonlReader, error) {
	records := []map[string]interface{}{}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(b, records)
	return &JsonlReader{records: records, startRow: startRow}, nil
}

type YamlReader struct {
	records  []map[string]interface{}
	f        *os.File
	counter  int
	startRow int
}

func (r *YamlReader) Read() ([]string, error) {
	if len(r.records) <= r.counter {
		return nil, io.EOF
	}
	if r.startRow > r.counter {
		r.counter++
		return nil, nil
	}
	record := r.records[r.counter]
	values := make([]string, len(record))
	r.counter++
	return values, nil
}

func (r *YamlReader) Close() error {
	return r.f.Close()
}

func newYamlReader(filename string, startRow int) (*YamlReader, error) {
	records := []map[string]interface{}{}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	yaml.Unmarshal(b, records)
	return &YamlReader{records: records, startRow: startRow}, nil
}

func getReader(c *cli.Context) (Reader, error) {
	f := c.String("file")
	encoding := c.String("encoding")
	start := c.Int("start-row")
	ext := filepath.Ext(f)
	mode := c.String("mode")

	var r Reader
	var err error
	switch ext {
	case ".csv", ".tsv":
		if ext == ".tsv" {
			mode = "tsv"
		}
		r, err = newCsvReader(f, encoding, mode, start)
	case ".xlsx":
		s := c.String("sheet")
		r, err = newExcelReader(f, s, start)
	case ".json":
		r, err = newJsonReader(f, start)
	case ".jsonl":
		r, err = newJsonlReader(f, start)
	case ".yaml", ".yml":
		r, err = newYamlReader(f, start)
	case ".dat":
		bs := strings.Split(c.String("bytes"), ",")
		bi := make([]int, len(bs))
		for i, b := range bs {
			bi[i], err = strconv.Atoi(b)
			if err != nil {
				return nil, err
			}
		}
		r, err = newFixWidthFileReader(f, encoding, bi)
	}
	return r, err
}
