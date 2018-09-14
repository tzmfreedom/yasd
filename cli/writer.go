package cli

import (
	"encoding/csv"
	"os"

	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/go-soapforce"
)

type writer interface {
	Write(record *soapforce.SObject) error
	Close()
}

type PPWriter struct {
	writer *csv.Writer
}

func (w *PPWriter) Write(record *soapforce.SObject) error {
	pp.Print(record)
	return nil
}

func (w *PPWriter) Close() {}

type CsvWriter struct {
	writer *csv.Writer
	fp     *os.File
}

func NewCsvWriter(cfg *config) (*CsvWriter, error) {
	var csvWriter *csv.Writer
	var writer *CsvWriter
	if cfg.Output != "" {
		fp, err := os.Create(cfg.Output)
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

func (w *CsvWriter) Close() {
	w.writer.Flush()
	if w.fp != nil {
		w.fp.Close()
	}
}

func getWriter(cfg *config) (writer, error) {
	switch cfg.Format {
	case "csv":
		return NewCsvWriter(cfg)
	case "tsv":
		return NewCsvWriter(cfg)
	case "debug":
		return &PPWriter{}, nil
	default:
		return NewCsvWriter(cfg)
	}
}
