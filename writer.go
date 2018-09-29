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
	Write(record *soapforce.SObject) error
	Close() error
}

type PPWriter struct {
	writer *csv.Writer
}

func (w *PPWriter) Write(record *soapforce.SObject) error {
	pp.Print(record)
	return nil
}

func (w *PPWriter) Close() error { return nil }

type CsvWriter struct {
	writer *csv.Writer
	fp     *os.File
}

func (w *CsvWriter) Write(record *soapforce.SObject) error {
	values := []string{}
	if record.Id != "" {
		values = append(values, record.Id)
	}
	for _, v := range record.Fields {
		values = append(values, v)
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

func newCsvWriter(f string, e string) (*CsvWriter, error) {
	var writer *CsvWriter
	if f != "" {
		fp, err := os.Create(f)
		if err != nil {
			return nil, err
		}
		w := newWriterWithEncoding(fp, e)
		csvWriter := csv.NewWriter(w)
		if runtime.GOOS == "windows" {
			csvWriter.UseCRLF = true
		}
		writer = &CsvWriter{
			writer: csvWriter,
			fp:     fp,
		}
	} else {
		w := newWriterWithEncoding(os.Stdout, e)
		csvWriter := csv.NewWriter(w)
		if runtime.GOOS == "windows" {
			csvWriter.UseCRLF = true
		}
		writer = &CsvWriter{
			writer: csvWriter,
		}
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
	f := c.String("output")
	format := c.String("format")
	e := c.String("encoding")

	switch format {
	case "csv", "tsv":
		return newCsvWriter(f, e)
	case "debug":
		return &PPWriter{}, nil
	default:
		return newCsvWriter(f, e)
	}
}
