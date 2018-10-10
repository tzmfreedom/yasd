package main

import (
	"testing"
)

func TestReadFromCsv(t *testing.T) {
	filename := "test/success.csv"
	encoding := "utf8"
	mode := ""
	start := 0
	reader, err := newCsvReader(filename, encoding, mode, start)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := []string{
		"あ",
		"i",
		"",
	}
	assertArrayEqual(t, reader, expected)

	expected = []string{
		"う",
		"e",
		"*",
	}
	assertArrayEqual(t, reader, expected)
	reader.Close()
}

func assertArrayEqual(t *testing.T, reader Reader, expected []string) {
	values, err := reader.Read()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	for i, actual := range values {
		if expected[i] != actual {
			t.Fatalf("expected '%s', but '%s'", expected[i], actual)
		}
	}
}

func TestReadFromExcel(t *testing.T) {
	filename := "test/success.xlsx"
	sheet := "test"
	start := 0
	reader, err := newExcelReader(filename, sheet, start)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := []string{
		"あ",
		"i",
		" ",
	}
	assertArrayEqual(t, reader, expected)

	expected = []string{
		"う",
		"e",
		"*",
	}
	assertArrayEqual(t, reader, expected)
	reader.Close()
}

func TestReadFromFixWidth(t *testing.T) {
	filename := "test/success.dat"
	encoding := "utf8"
	reader, err := newFixWidthFileReader(filename, encoding, []int{6, 3, 1})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := []string{
		"a",
		"bb",
		"",
	}
	assertArrayEqual(t, reader, expected)

	expected = []string{
		"あ",
		"い",
		"*",
	}
	assertArrayEqual(t, reader, expected)
	reader.Close()
}
