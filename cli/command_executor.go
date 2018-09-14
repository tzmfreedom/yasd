package cli

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tzmfreedom/go-soapforce"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"gopkg.in/yaml.v2"
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
	writer, err := getWriter(cfg)
	if err != nil {
		return err
	}
	defer writer.Close()
	for _, record := range res.Records {
		writer.Write(record)
	}
	for res.QueryLocator != "" {
		res, err := c.client.QueryMore(res.QueryLocator)
		if err != nil {
			return err
		}
		for _, record := range res.Records {
			writer.Write(record)
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

	sobjects := []*soapforce.SObject{}
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	headers, err = mapping(headers, cfg)
	if err != nil {
		return err
	}
	handler, err := getResponseHandler(cfg)
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
		sobject := createSObject(cfg.Type, headers, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			res, err := c.client.Create(sobjects)
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
	res, err := c.client.Create(sobjects)
	if err != nil {
		return  err
	}
	err = handler.Handle(res)
	return err
}

func (c *CommandExecutor) update(cfg *config) error {
	_, err := c.client.Login(cfg.Username, cfg.Password)
	if err != nil {
		return err
	}

	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := []*soapforce.SObject{}
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	headers, err = mapping(headers, cfg)
	if err != nil {
		return err
	}
	handler, err := getResponseHandler(cfg)
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
		sobject := createSObject(cfg.Type, headers, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			res, err := c.client.Update(sobjects)
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
	res, err := c.client.Update(sobjects)
	if err != nil {
		return err
	}
	err = handler.Handle(res)
	return err
}

func (c *CommandExecutor) upsert(cfg *config) error {
	_, err := c.client.Login(cfg.Username, cfg.Password)
	if err != nil {
		return err
	}

	reader, fp, err := getReader(cfg)
	if err != nil {
		return err
	}
	defer fp.Close()

	sobjects := []*soapforce.SObject{}
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	headers, err = mapping(headers, cfg)
	if err != nil {
		return err
	}
	handler, err := getResponseHandler(cfg)
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
		sobject := createSObject(cfg.Type, headers, fields)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			res, err := c.client.Upsert(sobjects, cfg.UpsertKey)
			if err != nil {
				return err
			}
			err = handler.HandleUpsert(res)
			if err != nil {
				return err
			}
			sobjects = sobjects[:0]
		}
	}
	res, err := c.client.Upsert(sobjects, cfg.UpsertKey)
	if err != nil {
		return err
	}
	err = handler.HandleUpsert(res)
	return err
}

func (c *CommandExecutor) delete(cfg *config) error {
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
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	headers, err = mapping(headers, cfg)
	if err != nil {
		return err
	}
	var ids []string
	handler, err := getResponseHandler(cfg)
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
			res, err := c.client.Delete(ids)
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
	res, err := c.client.Delete(ids)
	if err != nil {
		return err
	}
	err = handler.HandleDelete(res)
	return err
}

func (c *CommandExecutor) undelete(cfg *config) error {
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
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	headers, err = mapping(headers, cfg)
	if err != nil {
		return err
	}
	var ids []string
	handler, err := getResponseHandler(cfg)
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
			res, err := c.client.Undelete(ids)
			if err != nil {
				return nil
			}
			err = handler.HandleUndelete(res)
			if err != nil {
				return nil
			}
			ids = ids[:0]
		}
	}
	res, err := c.client.Undelete(ids)
	if err != nil {
		return err
	}
	err = handler.HandleUndelete(res)
	return err
}

func getReader(cfg *config) (*csv.Reader, *os.File, error) {
	fp, err := os.Open(cfg.InputFile)
	if err != nil {
		return nil, nil, err
	}
	var reader *csv.Reader
	switch strings.ToUpper(cfg.Encoding) {
	case "UTF8", "UTF-8":
		reader = csv.NewReader(fp)
	case "SHIFT-JIS", "SJIS":
		reader = csv.NewReader(transform.NewReader(fp, japanese.ShiftJIS.NewDecoder()))
	case "EUC-JP", "EUCJP":
		reader = csv.NewReader(transform.NewReader(fp, japanese.EUCJP.NewDecoder()))
	}
	reader.Comma = rune(cfg.Delimiter[0])
	reader.LazyQuotes = true
	return reader, fp, nil
}

func getId(headers []string, f []string) string {
	for i, header := range headers {
		if header == "Id" {
			return f[i]
		}
	}
	return ""
}

func createSObject(sObjectType string, headers []string, f []string) *soapforce.SObject {
	fields := map[string]string{}
	sobject := &soapforce.SObject{
		Type:   sObjectType,
		Fields: fields,
	}
	for i, header := range headers {
		if header == "Id" {
			sobject.Id = f[i]
		}
		fields[header] = f[i]
	}
	return sobject
}

func mapping(headers []string, cfg *config) ([]string, error) {
	if cfg.Mapping == "" {
		return nil, nil
	}
	mapping := map[string]string{}
	buf, err := ioutil.ReadFile(cfg.Mapping)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buf, mapping)
	if err != nil {
		return nil, err
	}
	returnHeaders := []string{}
	for _, h := range headers {
		if v, ok := mapping[h]; ok {
			returnHeaders = append(returnHeaders, v)
		} else {
			returnHeaders = append(returnHeaders, h)
		}
	}
	return returnHeaders, nil
}
