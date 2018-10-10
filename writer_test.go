package main

import (
	"testing"
	"github.com/tzmfreedom/go-soapforce"
	"bytes"
)

func TestCsvWrite(t *testing.T) {
	encoding := "utf8"
	comma := rune(',')
	buf := new(bytes.Buffer)
	writer, err := newCsvWriter(encoding, comma, buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	headers := []string{"あ","i",""}
	record := &soapforce.SObject{}
	err = writer.Write(headers, record)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestJsonWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	writer, err := newJsonWriter(buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	headers := []string{"あ","i",""}
	record := &soapforce.SObject{}
	err = writer.Write(headers, record)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestJsonlWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	writer, err := newJsonlWriter(buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	headers := []string{"あ","i",""}
	record := &soapforce.SObject{}
	err = writer.Write(headers, record)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

