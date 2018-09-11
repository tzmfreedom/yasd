package cli

import (
	"encoding/csv"
	"io"
	"os"

	"strings"

	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/go-soapforce"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type CommandExecutor struct {
	client *soapforce.Client
}

func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{
		client: soapforce.NewClient(),
	}
}

func (c *CommandExecutor) query(cfg *config) error {
	_, err := c.client.Login(cfg.Username, cfg.Password)
	if err != nil {
		return err
	}
	res, err := c.client.Query(cfg.Query)
	if err != nil {
		return err
	}
	writer, err := getWriter(cfg)
	if err != nil {
		return err
	}
	defer writer.Close()
	for _, record := range res.Records {
		writer.Write(record)
	}
	for res.QueryLocator != "" {
		res, err := c.client.QueryMore(res.QueryLocator)
		if err != nil {
			return err
		}
		for _, record := range res.Records {
			writer.Write(record)
		}
	}
	return nil
}

func (c *CommandExecutor) insert(cfg *config) error {
	_, err := c.client.Login(cfg.Username, cfg.Password)
	if err != nil {
		return err
	}

	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := []*soapforce.SObject{}
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createSObject(cfg.Type, headers, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			_, err = c.client.Create(sobjects)
			if err != nil {
				return err
			}
			sobjects = sobjects[:0]
		}
	}
	_, err = c.client.Create(sobjects)
	return err
}

func (c *CommandExecutor) update(cfg *config) error {
	_, err := c.client.Login(cfg.Username, cfg.Password)
	if err != nil {
		return err
	}

	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := []*soapforce.SObject{}
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createSObject(cfg.Type, headers, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			_, err = c.client.Create(sobjects)
			if err != nil {
				return err
			}
			sobjects = sobjects[:0]
		}
	}
	_, err = c.client.Create(sobjects)
	return err
}

func (c *CommandExecutor) upsert(cfg *config) error {
	_, err := c.client.Login(cfg.Username, cfg.Password)
	if err != nil {
		return err
	}

	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := []*soapforce.SObject{}
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createSObject(cfg.Type, headers, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			_, err = c.client.Upsert(sobjects, cfg.UpsertKey)
			if err != nil {
				return err
			}
			sobjects = sobjects[:0]
		}
	}
	_, err = c.client.Create(sobjects)
	return err
}

func (c *CommandExecutor) delete(cfg *config) error {
	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := make([]*soapforce.SObject, 0)
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createSObject(cfg.Type, fields, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			c.client.Create(sobjects)
		}
	}
	_, err = c.client.Create(sobjects)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) undelete(cfg *config) error {
	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := make([]*soapforce.SObject, 0)
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createSObject(cfg.Type, fields, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			c.client.Create(sobjects)
		}
	}
	_, err = c.client.Create(sobjects)
	if err != nil {
		return err
	}
	return nil
}

func getReader(cfg *config) (*csv.Reader, *os.File, error) {
	fp, err := os.Open(cfg.InputFile)
	if err != nil {
		return nil, nil, err
	}
	var reader *csv.Reader
	switch strings.ToUpper(cfg.Encoding) {
	case "UTF8", "UTF-8":
		reader = csv.NewReader(fp)
	case "SHIFT-JIS", "SJIS":
		reader = csv.NewReader(transform.NewReader(fp, japanese.ShiftJIS.NewDecoder()))
	case "EUC-JP", "EUCJP":
		reader = csv.NewReader(transform.NewReader(fp, japanese.EUCJP.NewDecoder()))
	}
	reader.Comma = rune(cfg.Delimiter[0])
	reader.LazyQuotes = true
	return reader, fp, nil
}

func createSObject(sObjectType string, headers []string, f []string) *soapforce.SObject {
	fields := map[string]string{}
	for i, header := range headers {
		fields[header] = f[i]
	}
	sobject := &soapforce.SObject{
		Type:   sObjectType,
		Fields: fields,
	}
	return sobject
}

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
			fp: fp,
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
