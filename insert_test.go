package main

import (
	"flag"
	"os"
	"testing"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/urfave/cli"
)

var username = os.Getenv("TEST_SALESFORCE_USERNAME")
var password = os.Getenv("TEST_SALESFORCE_PASSWORD")

func TestInsert(t *testing.T) {
	app := cli.NewApp()
	set := flag.NewFlagSet("test", flag.ExitOnError)
	for _, f := range insertFlags {
		f.Apply(set)
	}
	args := []string{
		"--username",
		username,
		"--password",
		password,
		"--file",
		"test/insert.csv",
		"--type",
		"Account",
	}
	err := set.Parse(args)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	ctx := cli.NewContext(app, set, nil)
	err = insert(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	assertRecords(t)
}

func assertRecords(t *testing.T) {
	client := soapforce.NewClient()
	_, err := client.Login(username, password)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	result, err := client.Query("SELECT Name FROM Account ORDER BY CreatedDate DESC LIMIT 1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(result.Records) != 1 {
		t.Fatalf("expected %d, but %d", 1, len(result.Records))
	}
	for _, record := range result.Records {
		actual := record.Fields["Name"]
		if actual != "test" {
			t.Fatalf("expected '%s', but '%s'", "test", actual)
		}
	}
}
