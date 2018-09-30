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

type responseHandler interface {
	Handle(results []*soapforce.SaveResult) error
	HandleUpsert(results []*soapforce.UpsertResult) error
	HandleDelete(results []*soapforce.DeleteResult) error
	HandleUndelete(results []*soapforce.UndeleteResult) error
}

type NoopResponseWriteHandler struct{}

func (h *NoopResponseWriteHandler) Handle(results []*soapforce.SaveResult) error         { return nil }
func (h *NoopResponseWriteHandler) HandleUpsert(results []*soapforce.UpsertResult) error { return nil }
func (h *NoopResponseWriteHandler) HandleDelete(results []*soapforce.DeleteResult) error { return nil }
func (h *NoopResponseWriteHandler) HandleUndelete(results []*soapforce.UndeleteResult) error {
	return nil
}

type ResponseWriteHandler struct {
	successWriter *csv.Writer
	errorWriter   *csv.Writer
}

func (h *ResponseWriteHandler) Handle(results []*soapforce.SaveResult) error {
	for _, result := range results {
		if result.Success {
			fields := []string{}
			fields = append(fields, result.Id)
			h.successWriter.Write(fields)
		} else {
			fields := []string{}
			errorMessages := []string{}
			for _, error := range result.Errors {
				errorMessages = append(errorMessages, error.Message)
			}
			pp.Print(result)
			errorMsg := strings.Join(errorMessages, ":")
			fields = append(fields, errorMsg)
			h.errorWriter.Write(fields)
		}
	}
	h.successWriter.Flush()
	h.errorWriter.Flush()
	return nil
}

func (h *ResponseWriteHandler) HandleUpsert(results []*soapforce.UpsertResult) error {
	for _, result := range results {
		if result.Success {
			fields := []string{}
			fields = append(fields, result.Id)
			h.successWriter.Write(fields)
		} else {
			fields := []string{}
			errorMessages := []string{}
			for _, error := range result.Errors {
				errorMessages = append(errorMessages, error.Message)
			}
			errorMsg := strings.Join(errorMessages, ":")
			fields = append(fields, errorMsg)
			h.errorWriter.Write(fields)
		}
	}
	h.successWriter.Flush()
	h.errorWriter.Flush()
	return nil
}

func (h *ResponseWriteHandler) HandleDelete(results []*soapforce.DeleteResult) error {
	for _, result := range results {
		if result.Success {
			fields := []string{}
			fields = append(fields, result.Id)
			h.successWriter.Write(fields)
		} else {
			fields := []string{}
			errorMessages := []string{}
			for _, error := range result.Errors {
				errorMessages = append(errorMessages, error.Message)
			}
			pp.Print(result)
			errorMsg := strings.Join(errorMessages, ":")
			fields = append(fields, errorMsg)
			h.errorWriter.Write(fields)
		}
	}
	h.successWriter.Flush()
	h.errorWriter.Flush()
	return nil
}

func (h *ResponseWriteHandler) HandleUndelete(results []*soapforce.UndeleteResult) error {
	for _, result := range results {
		if result.Success {
			fields := []string{}
			fields = append(fields, result.Id)
			h.successWriter.Write(fields)
		} else {
			fields := []string{}
			errorMessages := []string{}
			for _, error := range result.Errors {
				errorMessages = append(errorMessages, error.Message)
			}
			pp.Print(result)
			errorMsg := strings.Join(errorMessages, ":")
			fields = append(fields, errorMsg)
			h.errorWriter.Write(fields)
		}
	}
	h.successWriter.Flush()
	h.errorWriter.Flush()
	return nil
}

func newResponseWriteHandler(success string, error string, encoding string) (*ResponseWriteHandler, error) {
	successWriter, err := createCsvWriter(success, encoding)
	if err != nil {
		return nil, err
	}
	errorWriter, err := createCsvWriter(error, encoding)
	if err != nil {
		return nil, err
	}
	return &ResponseWriteHandler{
		successWriter: successWriter,
		errorWriter:   errorWriter,
	}, nil
}

func createCsvWriter(path string, encoding string) (*csv.Writer, error) {
	var writer *csv.Writer
	if path != "" {
		fp, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		var w io.Writer
		switch strings.ToUpper(encoding) {
		case "SHIFT-JIS", "SHIFT_JIS", "SJIS":
			w = transform.NewWriter(fp, japanese.ShiftJIS.NewEncoder())
		default:
			w = fp
		}
		writer = csv.NewWriter(w)
	} else {
		writer = csv.NewWriter(os.Stderr)
	}
	if runtime.GOOS == "windows" {
		writer.UseCRLF = true
	}
	return writer, nil
}

func getResponseHandler(c *cli.Context) (responseHandler, error) {
	success := c.String("success-file")
	error := c.String("error-file")
	encoding := c.String("encoding")

	h, err := newResponseWriteHandler(success, error, encoding)
	return h, err
}
