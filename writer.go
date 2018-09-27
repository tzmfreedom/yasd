package main

import (
	"encoding/csv"
	"os"

	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/go-soapforce"
	"github.com/urfave/cli"
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

func newCsvWriter(f string) (*CsvWriter, error) {
	var csvWriter *csv.Writer
	var writer *CsvWriter
	if f != "" {
		fp, err := os.Create(f)
		if err != nil {
			return nil, err
		}
		csvWriter = csv.NewWriter(fp)
		writer = &CsvWriter{
			writer: csvWriter,
			fp:     fp,
		}
	} else {
		csvWriter = csv.NewWriter(os.Stdout)
		writer = &CsvWriter{
			writer: csvWriter,
		}
	}
	return writer, nil
}

func getWriter(c *cli.Context) (writer, error) {
	f := c.String("output")
	format := c.String("format")
	switch format {
	case "csv":
		return newCsvWriter(f)
	case "tsv":
		return newCsvWriter(f)
	case "debug":
		return &PPWriter{}, nil
	default:
		return newCsvWriter(f)
	}
}
