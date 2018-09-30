package main

import (
	"github.com/urfave/cli"
)

var Commands = []cli.Command{
	{
		Name:    "export",
		Aliases: []string{"e"},
		Usage:   "Export SObject Record",
		Flags: append(
			defaultFlags(),
			cli.StringFlag{
				Name: "query, q",
			},
			cli.StringFlag{
				Name: "format",
			},
			cli.IntFlag{
				Name:  "batch-size",
				Value: 500,
			},
			cli.StringFlag{
				Name: "mode",
			},
		),
		Action: func(c *cli.Context) error {
			return query(c)
		},
	},
	{
		Name:    "insert",
		Aliases: []string{"i"},
		Usage:   "Insert SObject Record",
		Flags: append(
			defaultDmlFlags(),
			cli.BoolFlag{
				Name: "insert-nulls",
			},
		),
		Action: func(c *cli.Context) error {
			return insert(c)
		},
	},
	{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "Update SObject Record",
		Flags: append(
			defaultDmlFlags(),
			cli.BoolFlag{
				Name: "insert-nulls",
			},
		),
		Action: func(c *cli.Context) error {
			return update(c)
		},
	},
	{
		Name:  "upsert",
		Usage: "Upsert SObject Record",
		Flags: append(
			defaultDmlFlags(),
			cli.StringFlag{
				Name:  "upsert-key, k",
				Value: "Id",
			},
			cli.BoolFlag{
				Name: "insert-nulls",
			},
		),
		Action: func(c *cli.Context) error {
			return upsert(c)
		},
	},
	{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "Delete SObject Record",
		Flags:   defaultDmlFlags(),
		Action: func(c *cli.Context) error {
			return delete(c)
		},
	},
	{
		Name:  "undelete",
		Usage: "Undelete SObject Record",
		Flags: defaultDmlFlags(),
		Action: func(c *cli.Context) error {
			return undelete(c)
		},
	},
	{
		Name:  "generate-key",
		Usage: "Generate AES Key",
		Flags: defaultFlags(),
		Action: func(c *cli.Context) error {
			return generateEncryptionKey()
		},
	},
	{
		Name:  "encrypt",
		Usage: "Encrypt password",
		Flags: defaultFlags(),
		Action: func(c *cli.Context) error {
			return encryptCredential(c)
		},
	},
}

func defaultFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "username, u",
			EnvVar: "SALESFORCE_USERNAME",
		},
		cli.StringFlag{
			Name:   "password, p",
			EnvVar: "SALESFORCE_PASSWORD",
		},
		cli.StringFlag{
			Name:   "endpoint, e",
			Value:  "login.salesforce.com",
			EnvVar: "SALESFORCE_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "api-version",
			Value:  DefaultApiVersion,
			EnvVar: "SALESFORCE_APIVERSION",
		},
		cli.StringFlag{
			Name:  "mode",
			Value: "utf8",
		},
		cli.StringFlag{
			Name:  "encoding",
			Value: "utf8",
		},
		cli.StringFlag{
			Name: "mapping",
		},
		cli.BoolFlag{
			Name: "debug, d",
		},
		cli.StringFlag{
			Name: "config, c",
		},
		cli.StringFlag{
			Name: "key",
		},
	}
}

func defaultDmlFlags() []cli.Flag {
	return append(
		defaultFlags(),
		cli.StringFlag{
			Name: "file, f",
		},
		cli.StringFlag{
			Name: "type, t",
		},
		cli.StringFlag{
			Name:  "success-file",
			Value: "./success.csv",
		},
		cli.StringFlag{
			Name:  "error-file",
			Value: "./error.csv",
		},
	)
}
