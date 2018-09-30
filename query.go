package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/urfave/cli"
)

func query(c *cli.Context) error {
	if err := validateExportCommand(c); err != nil {
		return err
	}
	client := newClient(c)
	if err := login(client, c); err != nil {
		return err
	}

	q, err := buildQuery(client, c.String("query"))
	if err != nil {
		return err
	}
	res, err := client.Query(q)
	if err != nil {
		return err
	}
	writer, err := getWriter(c)
	if err != nil {
		return err
	}
	defer writer.Close()

	fields := getFields(q)
	writer.Header(fields)

	for _, record := range res.Records {
		writer.Write(fields, record)
	}
	for res.QueryLocator != "" {
		res, err := client.QueryMore(res.QueryLocator)
		if err != nil {
			return err
		}
		for _, record := range res.Records {
			writer.Write(fields, record)
		}
	}
	return nil
}

func buildQuery(c *soapforce.Client, original string) (string, error) {
	r := regexp.MustCompile(`(?i)SELECT\s+(\*)\s+FROM\s+([a-zA-Z\d_]+)`)
	matches := r.FindStringSubmatch(strings.TrimSpace(original))
	if len(matches) == 0 {
		return original, nil
	}
	result, err := c.DescribeSObject(matches[2])
	if err != nil {
		return "", err
	}
	fields := make([]string, len(result.Fields))
	for i, f := range result.Fields {
		fields[i] = f.Name
	}
	selectClause := strings.Join(fields, ",")
	return r.ReplaceAllString(original, fmt.Sprintf("SELECT %s FROM $2", selectClause)), nil
}

func getFields(q string) []string {
	r := regexp.MustCompile(`(?i)SELECT\s+([a-zA-Z\d_,\s]+)\sFROM\s`)
	matches := r.FindStringSubmatch(q)
	return strings.Split(matches[1], ",")
}

func validateExportCommand(c *cli.Context) error {
	if err := validateLoginFlag(c, "export"); err != nil {
		return err
	}
	q := c.String("query")
	if q == "" {
		_ = cli.ShowCommandHelp(c, "export")
		return cli.NewExitError("query is required", 1)
	}
	r := regexp.MustCompile(`(?i)SELECT\s+([a-zA-Z\d_,\s]+)\sFROM\s`)
	if !r.MatchString(q) {
		r := regexp.MustCompile(`(?i)SELECT\s+(\*)\s+FROM\s+([a-zA-Z\d_]+)`)
		if !r.MatchString(q) {
			_ = cli.ShowCommandHelp(c, "export")
			return cli.NewExitError("Malformed Query", 1)
		}
	}
	return nil
}
