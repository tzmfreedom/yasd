package main

import (
	"encoding/base64"
	"fmt"
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

func generateEncryptionKey() error {
	key, err := generateKey()
	if err != nil {
		return err
	}
	b64key := base64.StdEncoding.EncodeToString(key)
	fmt.Println(b64key)
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

func validateLoginFlag(c *cli.Context, command string) error {
	u := c.String("username")
	if u == "" {
		_ = cli.ShowCommandHelp(c, command)
		return cli.NewExitError("username is required", 1)
	}
	p := c.String("password")
	if p == "" {
		_ = cli.ShowCommandHelp(c, command)
		return cli.NewExitError("password is required", 1)
	}
	e := c.String("endpoint")
	if e == "" {
		_ = cli.ShowCommandHelp(c, command)
		return cli.NewExitError("endpoint is required", 1)
	}
	return nil
}
