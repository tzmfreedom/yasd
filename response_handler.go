package main

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/go-soapforce"
	"github.com/urfave/cli"
)

type responseHandler interface {
	Handle(results []*soapforce.SaveResult) error
	HandleUpsert(results []*soapforce.UpsertResult) error
	HandleDelete(results []*soapforce.DeleteResult) error
	HandleUndelete(results []*soapforce.UndeleteResult) error
}

type NoopResponseWriteHandler struct {}

func (h *NoopResponseWriteHandler) Handle(results []*soapforce.SaveResult) error { return nil }
func (h *NoopResponseWriteHandler) HandleUpsert(results []*soapforce.UpsertResult) error { return nil }
func (h *NoopResponseWriteHandler) HandleDelete(results []*soapforce.DeleteResult) error { return nil }
func (h *NoopResponseWriteHandler) HandleUndelete(results []*soapforce.UndeleteResult) error { return nil }

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

func newResponseWriteHandler(success string, error string) (*ResponseWriteHandler, error) {
	successWriter, err := createCsvWriter(success)
	if err != nil {
		return nil, err
	}
	errorWriter, err := createCsvWriter(error)
	if err != nil {
		return nil, err
	}
	return &ResponseWriteHandler{
		successWriter: successWriter,
		errorWriter:   errorWriter,
	}, nil
}

func createCsvWriter(path string) (*csv.Writer, error) {
	var writer *csv.Writer
	if path != "" {
		fp, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		writer = csv.NewWriter(fp)
	} else {
		writer = csv.NewWriter(os.Stderr)
	}
	return writer, nil
}

func getResponseHandler(c *cli.Context) (responseHandler, error) {
	success := c.String("success-file")
	error := c.String("error-file")
	h, err := newResponseWriteHandler(success, error)
	return h, err
}