package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func newClient(c *cli.Context) *soapforce.Client {
	client := soapforce.NewClient()
	client.SetDebug(c.Bool("debug"))
	client.SetLoginUrl(c.String("endpoint"))
	client.SetBatchSize(c.Int("batch-size"))
	return client
}

func login(client *soapforce.Client, ctx *cli.Context) error {
	var err error
	username := ctx.String("username")
	password := ctx.String("password")
	keypath := ctx.String("key")
	if keypath != "" {
		password, err = decryptCredential(keypath, password)
		if err != nil {
			return err
		}
	}
	_, err = client.Login(username, password)
	return err
}

func query(c *cli.Context) error {
	client := newClient(c)
	if err := login(client, c); err != nil {
		return err
	}

	q := c.String("query")
	res, err := client.Query(q)
	if err != nil {
		return err
	}
	writer, err := getWriter(c)
	if err != nil {
		return err
	}
	defer writer.Close()
	for _, record := range res.Records {
		writer.Write(record)
	}
	for res.QueryLocator != "" {
		res, err := client.QueryMore(res.QueryLocator)
		if err != nil {
			return err
		}
		for _, record := range res.Records {
			writer.Write(record)
		}
	}
	return nil
}

func insert(c *cli.Context) error {
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
		sobject := createInsertSObject(t, headers, fields, insertNulls)
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

func update(c *cli.Context) error {
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
		sobject := createSObject(t, headers, fields, insertNulls)
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

func upsert(c *cli.Context) error {
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
	upsertKey := c.String("upsert-key")
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sobject := createSObject(t, headers, fields, insertNulls)
		sobjects = append(sobjects, sobject)
		if len(sobjects) == 200 {
			res, err := client.Upsert(sobjects, upsertKey)
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
	res, err := client.Upsert(sobjects, upsertKey)
	if err != nil {
		return err
	}
	err = handler.HandleUpsert(res)
	return err
}

func delete(c *cli.Context) error {
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

func undelete(c *cli.Context) error {
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
			res, err := client.Undelete(ids)
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
	res, err := client.Undelete(ids)
	if err != nil {
		return err
	}
	err = handler.HandleUndelete(res)
	return err
}

func generateEncryptionKey() error {
	key, err := generateKey()
	if err != nil {
		return err
	}
	b64key := base64.StdEncoding.EncodeToString(key)
	fmt.Println(b64key)
	return nil
}

func encryptCredential(c *cli.Context) error {
	keypath := c.String("key")
	password := c.String("password")
	b64key, err := ioutil.ReadFile(keypath)
	if err != nil {
		return err
	}
	key, err := base64.StdEncoding.DecodeString(string(b64key))
	if err != nil {
		return err
	}
	encryptedPassword, err := encrypt([]byte(password), key)
	if err != nil {
		return err
	}
	fmt.Println(encryptedPassword)
	return nil
}

func decryptCredential(keypath string, password string) (string, error) {
	b64key, err := ioutil.ReadFile(keypath)
	if err != nil {
		return "", err
	}
	key, err := base64.StdEncoding.DecodeString(string(b64key))
	if err != nil {
		return "", err
	}
	return decrypt(password, key)
}

func getId(headers []string, f []string) string {
	for i, header := range headers {
		if header == "Id" {
			return f[i]
		}
	}
	return ""
}

func createSObject(sObjectType string, headers []string, f []string, insertNulls bool) *soapforce.SObject {
	fields := map[string]string{}
	sobject := &soapforce.SObject{
		Type:   sObjectType,
		Fields: fields,
	}
	fieldsToNull := []string{}
	for i, header := range headers {
		if header == "Id" {
			sobject.Id = f[i]
		} else if insertNulls && f[i] == "" {
			fieldsToNull = append(fieldsToNull, header)
		} else {
			fields[header] = f[i]
		}
	}
	sobject.FieldsToNull = fieldsToNull
	return sobject
}

func createInsertSObject(sObjectType string, headers []string, f []string, insertNulls bool) *soapforce.SObject {
	fields := map[string]string{}
	sobject := &soapforce.SObject{
		Type:   sObjectType,
		Fields: fields,
	}
	fieldsToNull := []string{}
	for i, header := range headers {
		if header != "Id" {
			if insertNulls && f[i] == "" {
				fieldsToNull = append(fieldsToNull, header)
			} else {
				fields[header] = f[i]
			}
		}
	}
	sobject.FieldsToNull = fieldsToNull
	return sobject
}

func mapping(headers []string, m string) ([]string, error) {
	if m == "" {
		return headers, nil
	}
	mapping := map[string]string{}
	buf, err := ioutil.ReadFile(m)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(buf, mapping); err != nil {
		return nil, err
	}
	returnHeaders := make([]string, len(headers))
	for i, h := range headers {
		if v, ok := mapping[h]; ok {
			returnHeaders[i] = v
		} else {
			returnHeaders[i] = h
		}
	}
	return returnHeaders, nil
}
