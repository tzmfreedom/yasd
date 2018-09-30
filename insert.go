package main

import (
	"io"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/go-soapforce"
	"github.com/urfave/cli"
)

func insert(c *cli.Context) error {
	if err := validateInsertCommand(c); err != nil {
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
	setReferenceMap(client, t)

	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createInsertSObject(client, t, headers, fields, insertNulls)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			res, err := client.Create(sobjects)
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
	res, err := client.Create(sobjects)
	if err != nil {
		return err
	}
	err = handler.Handle(res)
	return err
}

func validateInsertCommand(c *cli.Context) error {
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

func createInsertSObject(client *soapforce.Client, sObjectType string, headers []string, f []string, insertNulls bool) *soapforce.SObject {
	fields := map[string]interface{}{}
	sobject := &soapforce.SObject{
		Type:   sObjectType,
		Fields: fields,
	}
	fieldsToNull := []string{}
	for i, header := range headers {
		if header != "Id" {
			if strings.Contains(header, ".") {
				values := strings.Split(header, ".")
				referenceField := strings.Replace(values[0], "__R", "__r", -1)

				obj := map[string]string{}
				pp.Print(globalReferenceMap)
				obj["type"] = globalReferenceMap[client.UserInfo.OrganizationId][referenceField]
				obj[values[1]] = f[i]
				fields[values[0]] = obj
			} else {
				if insertNulls && f[i] == "" {
					fieldsToNull = append(fieldsToNull, header)
				} else {
					fields[header] = f[i]
				}
			}
		}
	}
	sobject.FieldsToNull = fieldsToNull
	return sobject
}
