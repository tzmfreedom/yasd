package main

import (
	"encoding/csv"
	"io"
	"os"
	"runtime"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"encoding/json"

	"github.com/k0kubun/pp"
	"github.com/tealeg/xlsx"
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
	m := newCaseInsensitiveMap(record.Fields)
	values := make([]string, len(headers))
	for i, h := range headers {
		if strings.ToLower(h) == "id" {
			values[i] = record.Id
		} else {
			if strings.Contains(h, ".") {
				values[i] = getField(m, h)
			} else {
				pp.Println(h)
				pp.Println(m.Get(h))
				pp.Println(m)
				values[i] = m.Get(h).(string)
			}
		}
	}
	return w.writer.Write(values)
}

func getField(m *caseInsensitiveMap, h string) string {
	keys := strings.Split(h, ".")
	fields := keys[:len(keys)-2]
	v := m.Get(keys[0])
	if v == nil {
		return ""
	}
	sobj := v.(*soapforce.SObject)

	for _, field := range fields {
		m = newCaseInsensitiveMap(sobj.Fields)
		v = m.Get(field)
		if v == nil {
			return ""
		}
		sobj = v.(*soapforce.SObject)
	}

	m = newCaseInsensitiveMap(sobj.Fields)
	v = m.Get(keys[len(keys)-1])
	if v == nil {
		return ""
	}
	return v.(string)
}

func (w *CsvWriter) Close() error {
	w.writer.Flush()
	if w.fp != nil {
		return w.fp.Close()
	}
	return nil
}

func newCsvWriter(e string, comma rune, ioWriter io.Writer) (*CsvWriter, error) {
	w := newWriterWithEncoding(ioWriter, e)
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

type caseInsensitiveMap struct {
	fields map[string]interface{}
}

func (m *caseInsensitiveMap) Get(k string) interface{} {
	return m.fields[strings.ToLower(k)]
}

func (m *caseInsensitiveMap) Set(k string, v interface{}) {
	m.fields[strings.ToLower(k)] = v
}

func newCaseInsensitiveMap(m map[string]interface{}) *caseInsensitiveMap {
	cim := &caseInsensitiveMap{fields: map[string]interface{}{}}
	for k, v := range m {
		cim.Set(k, v)
	}
	return cim
}

type JsonWriter struct {
	e *json.Encoder
	r []map[string]interface{}
}

func newJsonWriter(writer io.Writer) (*JsonWriter, error) {
	e := json.NewEncoder(writer)
	r := []map[string]interface{}{}
	return &JsonWriter{e: e, r: r}, nil
}

func (w *JsonWriter) Header(h []string) error {
	return nil
}

func (w *JsonWriter) Write(headers []string, record *soapforce.SObject) error {
	f := getWriteFields(record)
	w.r = append(w.r, f)
	return nil
}

func (w *JsonWriter) Close() error {
	return w.e.Encode(w.r)
}

type JsonlWriter struct {
	e *json.Encoder
}

func newJsonlWriter(writer io.Writer) (*JsonlWriter, error) {
	e := json.NewEncoder(writer)
	return &JsonlWriter{e: e}, nil
}

func (w *JsonlWriter) Header(h []string) error {
	return nil
}

func (w *JsonlWriter) Write(headers []string, record *soapforce.SObject) error {
	f := getWriteFields(record)
	return w.e.Encode(f)
}

func (w *JsonlWriter) Close() error {
	return nil
}

type YamlWriter struct {
	e *yaml.Encoder
	r []map[string]interface{}
}

func newYamlWriter(writer io.Writer) (*YamlWriter, error) {
	e := yaml.NewEncoder(writer)
	r := []map[string]interface{}{}
	return &YamlWriter{e: e, r: r}, nil
}

func (w *YamlWriter) Header(h []string) error {
	return nil
}

func (w *YamlWriter) Write(headers []string, record *soapforce.SObject) error {
	f := getWriteFields(record)
	w.r = append(w.r, f)
	return nil
}

func (w *YamlWriter) Close() error {
	return w.e.Encode(w.r)
}

type XlsxWriter struct {
	s     *xlsx.Sheet
	f     *xlsx.File
	fName string
}

func (w *XlsxWriter) Header(headers []string) error {
	row := w.s.AddRow()
	for _, h := range headers {
		cell := row.AddCell()
		cell.Value = h
	}
	return nil
}

func (w *XlsxWriter) Write(headers []string, record *soapforce.SObject) error {
	m := newCaseInsensitiveMap(record.Fields)
	row := w.s.AddRow()
	for _, h := range headers {
		cell := row.AddCell()
		if strings.ToLower(h) == "id" {
			cell.Value = record.Id
		} else {
			if strings.Contains(h, ".") {
				cell.Value = getField(m, h)
			} else {
				cell.Value = m.Get(h).(string)
			}
		}
	}
	return nil
}

func (w *XlsxWriter) Close() error {
	return w.f.Save(w.fName)
}

func newXlsxWriter(fName string, sheetName string) (*XlsxWriter, error) {
	var f *xlsx.File
	if _, err := os.Stat(fName); err != nil {
		f = xlsx.NewFile()
	} else {
		f, err = xlsx.OpenFile(fName)
		if err != nil {
			return nil, err
		}
	}
	s, err := f.AddSheet(sheetName)
	if err != nil {
		return nil, err
	}
	return &XlsxWriter{f: f, s: s, fName: fName}, nil
}

func getWriteFields(r *soapforce.SObject) map[string]interface{} {
	f := map[string]interface{}{}
	if r.Id != "" {
		f["Id"] = r.Id
	}
	for k, v := range r.Fields {
		if sv, ok := v.(string); ok {
			f[k] = sv
		} else if sobj, ok := v.(*soapforce.SObject); ok {
			f[k] = getWriteFields(sobj)
		} else {
			f[k] = v
		}
	}
	return f
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
	switch format {
	case "tsv", "t":
		comma = '\t'
	default:
		comma = ','
	}

	switch format {
	case "csv", "tsv":
		return newCsvWriter(e, comma, os.Stdout)
	case "jsonl":
		return newJsonlWriter(os.Stdout)
	case "json":
		return newJsonWriter(os.Stdout)
	case "yaml", "yml":
		return newYamlWriter(os.Stdout)
	case "xlsx":
		fName := c.String("file")
		s := c.String("sheet")
		return newXlsxWriter(fName, s)
	case "debug":
		return &PPWriter{}, nil
	default:
		return newCsvWriter(e, comma, os.Stdout)
	}
}
