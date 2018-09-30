package main

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"bufio"

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

type ExcelReader struct {
	xf       *xlsx.File
	xs       *xlsx.Sheet
	counter  int
	startRow int
}

func (r *ExcelReader) Read() ([]string, error) {
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

func newFixWidthFileReader(f string, e string) (*FixWidthFileReader, error) {
	fp, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(fp)
	return &FixWidthFileReader{f: fp, s: s, e: e, byteNumbers: []int{6, 7, 3}}, nil
}

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
			}, nil
		}
	}
	return nil, errors.New("Sheet does not exists")
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
	case ".csv":
		fp, err := os.Open(f)
		if err != nil {
			return nil, err
		}
		var cr *csv.Reader
		switch strings.ToUpper(encoding) {
		case "UTF8", "UTF-8":
			cr = csv.NewReader(fp)
		case "SHIFT-JIS", "SJIS":
			cr = csv.NewReader(transform.NewReader(fp, japanese.ShiftJIS.NewDecoder()))
		case "EUC-JP", "EUCJP":
			cr = csv.NewReader(transform.NewReader(fp, japanese.EUCJP.NewDecoder()))
		}
		if mode == "tsv" {
			cr.Comma = '\t'
		}
		cr.LazyQuotes = true

		r = &CsvReader{cr: cr, f: fp, startRow: start}
	case ".xls":
		s := "import"
		r, err = newExcelReader(f, s, start)
		if err != nil {
			return nil, err
		}
	case ".xlsx":
	case ".json":
	}
	return r, nil
}
