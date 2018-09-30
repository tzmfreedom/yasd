package main

import (
	"encoding/csv"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/go-soapforce"
	"github.com/urfave/cli"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type writer interface {
	Header([]string) error
	Write(headers []string, record *soapforce.SObject) error
	Close() error
}

type PPWriter struct {
	writer *csv.Writer
}

func (w *PPWriter) Header(h []string) error {
	pp.Print(h)
	return nil
}

func (w *PPWriter) Write(headers []string, record *soapforce.SObject) error {
	pp.Print(record)
	return nil
}

func (w *PPWriter) Close() error { return nil }

type CsvWriter struct {
	writer *csv.Writer
	fp     *os.File
}

func (w *CsvWriter) Header(h []string) error {
	return w.writer.Write(h)
}

func (w *CsvWriter) Write(headers []string, record *soapforce.SObject) error {
	values := make([]string, len(headers))
	for i, h := range headers {
		if h == "Id" {
			values[i] = record.Id
		} else {
			values[i] = record.Fields[h].(string)
		}
	}
	return w.writer.Write(values)
}

func (w *CsvWriter) Close() error {
	w.writer.Flush()
	if w.fp != nil {
		return w.fp.Close()
	}
	return nil
}

func newCsvWriter(e string, comma rune) (*CsvWriter, error) {
	w := newWriterWithEncoding(os.Stdout, e)
	csvWriter := csv.NewWriter(w)
	if runtime.GOOS == "windows" {
		csvWriter.UseCRLF = true
	}
	csvWriter.Comma = comma
	writer := &CsvWriter{
		writer: csvWriter,
	}
	return writer, nil
}

func newWriterWithEncoding(w io.Writer, e string) io.Writer {
	switch strings.ToUpper(e) {
	case "SHIFT-JIS", "SHIFT_JIS", "SJIS":
		return transform.NewWriter(w, japanese.ShiftJIS.NewEncoder())
	default:
		return w
	}
}

func getWriter(c *cli.Context) (writer, error) {
	format := c.String("format")
	e := c.String("encoding")
	var comma rune
	switch c.String("mode") {
	case "tsv", "t":
		comma = '\t'
	default:
		comma = ','
	}

	switch format {
	case "csv", "tsv":
		return newCsvWriter(e, comma)
	case "debug":
		return &PPWriter{}, nil
	default:
		return newCsvWriter(e, comma)
	}
}
