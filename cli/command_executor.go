package cli

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/go-soapforce"
)

type CommandExecutor struct {
	client *soapforce.Client
}

func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{
		client: soapforce.NewClient(),
	}
}

func (c *CommandExecutor) query(cfg *config) error {
	_, err := c.client.Login(cfg.Username, cfg.Password)
	if err != nil {
		return err
	}
	res, err := c.client.Query(cfg.Query)
	if err != nil {
		return err
	}
	for _, record := range res.Records {
		outputRecord(record, "csv")
	}
	for res.QueryLocator != "" {
		res, err := c.client.QueryMore(res.QueryLocator)
		if err != nil {
			return err
		}
		for _, record := range res.Records {
			outputRecord(record, "csv")
		}
	}
	return nil
}

func (c *CommandExecutor) insert(cfg *config) error {
	_, err := c.client.Login(cfg.Username, cfg.Password)
	if err != nil {
		return err
	}

	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := make([]*soapforce.SObject, 0)
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createSObject(cfg.Type, fields, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			c.client.Create(sobjects)
		}
	}
	_, err = c.client.Create(sobjects)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) update(cfg *config) error {
	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := make([]*soapforce.SObject, 0)
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createSObject(cfg.Type, fields, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			c.client.Create(sobjects)
		}
	}
	_, err = c.client.Create(sobjects)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) upsert(cfg *config) error {
	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := make([]*soapforce.SObject, 0)
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createSObject(cfg.Type, fields, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			c.client.Create(sobjects)
		}
	}
	_, err = c.client.Create(sobjects)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) delete(cfg *config) error {
	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := make([]*soapforce.SObject, 0)
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createSObject(cfg.Type, fields, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			c.client.Create(sobjects)
		}
	}
	_, err = c.client.Create(sobjects)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) undelete(cfg *config) error {
	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := make([]*soapforce.SObject, 0)
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createSObject(cfg.Type, fields, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			c.client.Create(sobjects)
		}
	}
	_, err = c.client.Create(sobjects)
	if err != nil {
		return err
	}
	return nil
}

func getReader(cfg *config) (*csv.Reader, *os.File, error) {
	fp, err := os.Open(cfg.InputFile)
	if err != nil {
		return nil, nil, err
	}
	reader := csv.NewReader(fp)
	reader.Comma = rune(cfg.Delimiter[0])
	reader.LazyQuotes = true
	return reader, fp, nil
}

func createSObject(sObjectType string, headers []string, f []string) *soapforce.SObject {
	fields := map[string]string{}
	for i, header := range headers {
		fields[header] = f[i]
	}
	sobject := &soapforce.SObject{
		Type:   sObjectType,
		Fields: fields,
	}
	return sobject
}

func outputRecord(record *soapforce.SObject, formatter string) {
	pp.Print(record)
}
