package main

import (
	"io"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/urfave/cli"
)

func update(c *cli.Context) error {
	if err := validateUpdateCommand(c); err != nil {
		return err
	}
	client := newClient(c)
	if err := login(client, c); err != nil {
		return err
	}

	reader, err := getReader(c)
	if err != nil {
		return err
	}
	defer reader.Close()

	sobjects := []*soapforce.SObject{}
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	headers, err = mapping(headers, c.String("mapping"))
	if err != nil {
		return err
	}
	handler, err := getResponseHandler(c)
	if err != nil {
		return err
	}
	t := c.String("type")
	insertNulls := c.Bool("insert-nulls")
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if fields == nil {
			continue
		}
		sobject := createSObject(client, t, headers, fields, insertNulls)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			res, err := client.Update(sobjects)
			if err != nil {
				return err
			}
			err = handler.Handle(res)
			if err != nil {
				return err
			}
			sobjects = sobjects[:0]
		}
	}
	res, err := client.Update(sobjects)
	if err != nil {
		return err
	}
	err = handler.Handle(res)
	return err
}

func validateUpdateCommand(c *cli.Context) error {
	if err := validateLoginFlag(c, "insert"); err != nil {
		return err
	}
	t := c.String("type")
	if t == "" {
		_ = cli.ShowCommandHelp(c, "insert")
		return cli.NewExitError("type is required", 1)
	}
	f := c.String("file")
	if f == "" {
		_ = cli.ShowCommandHelp(c, "insert")
		return cli.NewExitError("file is required", 1)
	}
	return nil
}
