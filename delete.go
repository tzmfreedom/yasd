package main

import (
	"io"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/urfave/cli"
)

func delete(c *cli.Context) error {
	if err := validateDeleteCommand(c); err != nil {
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

	sobjects := make([]*soapforce.SObject, 0)
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	headers, err = mapping(headers, c.String("mapping"))
	if err != nil {
		return err
	}
	var ids []string
	handler, err := getResponseHandler(c)
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
		id := getId(headers, fields)
		ids = append(ids, id)
		if len(sobjects) == 200 {
			res, err := client.Delete(ids)
			if err != nil {
				return err
			}
			err = handler.HandleDelete(res)
			if err != nil {
				return err
			}
			ids = ids[:0]
		}
	}
	res, err := client.Delete(ids)
	if err != nil {
		return err
	}
	err = handler.HandleDelete(res)
	return err
}

func validateDeleteCommand(c *cli.Context) error {
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
