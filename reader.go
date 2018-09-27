package main

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
	"github.com/urfave/cli"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Reader interface {
	Read() ([]string, error)
	Close() error
}

type CsvReader struct {
	cr *csv.Reader
	f *os.File
}

func (r *CsvReader) Read() ([]string, error) {
	return r.Read()
}

func (r *CsvReader) Close() error {
	return r.f.Close()
}

type ExcelReader struct {
	xf      *xlsx.File
	xs      *xlsx.Sheet
	counter int
}

func (r *ExcelReader) Read() ([]string, error) {
	row := r.xs.Rows[r.counter]
	values := make([]string, len(row.Cells))
	for i, cell := range row.Cells {
		values[i] = cell.Value
	}
	r.counter++
	return values, nil
}

func (r *ExcelReader) Close() error { return nil }

func newExcelReader(f string, sheet string) (*ExcelReader, error) {
	xf, err := xlsx.OpenFile(f)
	if err != nil {
		return nil, err
	}
	for _, s := range xf.Sheets {
		if s.Name == sheet {
			return &ExcelReader{
				counter: 0,
				xf:      xf,
				xs:      s,
			}, nil
		}
	}
	return nil, errors.New("Sheet does not exists")
}

func getReader(c *cli.Context) (Reader, error) {
	f := c.String("file")
	encoding := c.String("encoding")
	delimiter := c.String("delimiter")[0]
	filetype := c.String("file-type")

	var r Reader
	var err error
	switch filetype {
	case "csv":
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
		cr.Comma = rune(delimiter)
		cr.LazyQuotes = true

		r = &CsvReader{ cr: cr, f: fp }
	case "xls":
		s := "import"
		r, err = newExcelReader(f, s)
		if err != nil {
			return nil, err
		}
	case "xlsx":
	case "json":
	}
	return r, nil
}
