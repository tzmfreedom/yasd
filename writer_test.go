package main

import (
	"bytes"
	"testing"

	"github.com/tzmfreedom/go-soapforce"
)

func TestCsvHeader(t *testing.T) {
	encoding := "utf8"
	comma := rune(',')
	buf := new(bytes.Buffer)
	writer, err := newCsvWriter(encoding, comma, buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	headers := []string{"あ", "i", ""}
	err = writer.Header(headers)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := "あ,i,\n"
	if buf.String() != expected {
		t.Fatalf("expected: '%s', but '%s'", expected, buf.String())
	}
}

func TestCsvWrite(t *testing.T) {
	encoding := "utf8"
	comma := rune(',')
	buf := new(bytes.Buffer)
	writer, err := newCsvWriter(encoding, comma, buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	headers := []string{"Name", "A__c"}
	record := &soapforce.SObject{
		Fields: map[string]interface{}{
			"Name": "aaa",
			"A__c": ",,",
		},
	}
	err = writer.Write(headers, record)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := "aaa,\",,\"\n"
	if buf.String() != expected {
		t.Fatalf("expected: '%s', but '%s'", expected, buf.String())
	}
}

func TestJsonWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	writer, err := newJsonWriter(buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	headers := []string{"あ", "i", ""}
	record := &soapforce.SObject{
		Fields: map[string]interface{}{
			"Name": "aaa",
			"A__c": ",,",
		},
	}
	err = writer.Write(headers, record)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := "[{\"A__c\":\",,\",\"Name\":\"aaa\"}]\n"
	if buf.String() != expected {
		t.Fatalf("expected: '%s', but '%s'", expected, buf.String())
	}
}

func TestJsonlWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	writer, err := newJsonlWriter(buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	headers := []string{"あ", "i", ""}
	record := &soapforce.SObject{
		Fields: map[string]interface{}{
			"Name": "aaa",
			"A__c": ",,",
		},
	}
	err = writer.Write(headers, record)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := "{\"A__c\":\",,\",\"Name\":\"aaa\"}\n"
	if buf.String() != expected {
		t.Fatalf("expected: '%s', but '%s'", expected, buf.String())
	}
}

func TestYamlWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	writer, err := newYamlWriter(buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	headers := []string{"あ", "i", ""}
	record := &soapforce.SObject{
		Fields: map[string]interface{}{
			"Name": "aaa",
			"A__c": ",,",
		},
	}
	err = writer.Write(headers, record)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := "- A__c: ',,'\n  Name: aaa\n"
	if buf.String() != expected {
		t.Fatalf("expected: '%s', but '%s'", expected, buf.String())
	}
}

func TestCaseInsentivieGet(t *testing.T) {
	m := map[string]interface{}{
		"AAA": "aaa",
		"BBB": "bbb",
		"CcC": "ccc",
		"ddd": "DDD",
	}
	cim := newCaseInsensitiveMap(m)
	for _, v := range m {
		actual := cim.Get(v.(string))
		expected := v.(string)
		if actual != expected {
			t.Fatalf("expected: '%s', but '%s'", expected, actual)
		}
	}
}

func TestCaseInsentivieSet(t *testing.T) {
	m := map[string]interface{}{}
	cim := newCaseInsensitiveMap(m)
	cim.Set("aaa", "bbb")

	actual := cim.Get("AAA")
	expected := "bbb"
	if actual != expected {
		t.Fatalf("expected: '%s', but '%s'", expected, actual)
	}
	actual = cim.Get("aaa")
	if actual != expected {
		t.Fatalf("expected: '%s', but '%s'", expected, actual)
	}
}
