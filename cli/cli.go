package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

type CLI struct {
	Config *config
}

type config struct {
	Username          string
	Password          string
	Endpoint          string
	ApiVersion        string
	Query             string
	Type              string
	InputFile         string
	Delimiter         string
	Encoding          string
	UpsertKey         string
	Output            string
	Format            string
	Mapping           string
	ErrorPath         string
	SuccessPath       string
	EncryptionKeyPath string
	BatchSize         int
	Debug             bool
	InsertNulls       bool
	UpdateKey         bool
	ConfigFile        string
}

var (
	Version  string
	Revision string
)

const (
	AppName = "yasd"
	Usage   = "Yet Another Salesforce Dataloader"
)

const (
	DefaultApiVersion = "38.0"
)

func NewCli() *CLI {
	c := &CLI{
		Config: &config{},
	}
	return c
}

func (c *CLI) Run(args []string) error {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s revision=%s\n", c.App.Version, Revision)
	}

	app := cli.NewApp()
	app.Name = AppName
	app.Usage = Usage
	app.Version = Version

	app.Commands = []cli.Command{
		{
			Name:    "export",
			Aliases: []string{"e"},
			Usage:   "Export SObject Record",
			Flags: append(
				c.defaultFlags(),
				cli.StringFlag{
					Name:        "query, q",
					Destination: &c.Config.Query,
				},
				cli.StringFlag{
					Name:        "output, o",
					Destination: &c.Config.Output,
				},
				cli.StringFlag{
					Name:        "format",
					Destination: &c.Config.Format,
				},
				cli.IntFlag{
					Name:        "batch-size",
					Value:       500,
					Destination: &c.Config.BatchSize,
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.query(*c.Config)
			},
		},
		{
			Name:    "insert",
			Aliases: []string{"i"},
			Usage:   "Insert SObject Record",
			Flags: append(
				c.defaultDmlFlags(),
				cli.BoolFlag{
					Name:        "insert-nulls",
					Destination: &c.Config.InsertNulls,
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.insert(*c.Config)
			},
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update SObject Record",
			Flags: append(
				c.defaultDmlFlags(),
				cli.BoolFlag{
					Name:        "insert-nulls",
					Destination: &c.Config.InsertNulls,
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.update(*c.Config)
			},
		},
		{
			Name:  "upsert",
			Usage: "Upsert SObject Record",
			Flags: append(
				c.defaultDmlFlags(),
				cli.StringFlag{
					Name:        "upsert-key, k",
					Value:       "Id",
					Destination: &c.Config.UpsertKey,
				},
				cli.BoolFlag{
					Name:        "insert-nulls",
					Destination: &c.Config.InsertNulls,
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.upsert(*c.Config)
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete SObject Record",
			Flags:   c.defaultDmlFlags(),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.delete(*c.Config)
			},
		},
		{
			Name:  "undelete",
			Usage: "Undelete SObject Record",
			Flags: c.defaultDmlFlags(),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.undelete(*c.Config)
			},
		},
		{
			Name:  "generate-key",
			Usage: "Generate AES Key",
			Flags: append(
				c.defaultFlags(),
				cli.StringFlag{
					Name:        "key",
					Destination: &c.Config.EncryptionKeyPath,
				},
				cli.BoolFlag{
					Name:        "force-update",
					Destination: &c.Config.UpdateKey,
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.generateEncryptionKey(*c.Config)
			},
		},
		{
			Name:  "encrypt",
			Usage: "Encrypt password",
			Flags: c.defaultFlags(),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.encryptCredential(*c.Config)
			},
		},
	}
	return app.Run(args)
}

func (c *CLI) defaultFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "username, u",
			Destination: &c.Config.Username,
			EnvVar:      "SALESFORCE_USERNAME",
		},
		cli.StringFlag{
			Name:        "password, p",
			Destination: &c.Config.Password,
			EnvVar:      "SALESFORCE_PASSWORD",
		},
		cli.StringFlag{
			Name:        "endpoint, e",
			Value:       "login.salesforce.com",
			Destination: &c.Config.Endpoint,
			EnvVar:      "SALESFORCE_ENDPOINT",
		},
		cli.StringFlag{
			Name:        "api-version",
			Value:       DefaultApiVersion,
			Destination: &c.Config.ApiVersion,
			EnvVar:      "SALESFORCE_APIVERSION",
		},
		cli.StringFlag{
			Name:        "delimiter",
			Value:       ",",
			Destination: &c.Config.Delimiter,
		},
		cli.StringFlag{
			Name:        "encoding",
			Value:       "utf8",
			Destination: &c.Config.Encoding,
		},
		cli.StringFlag{
			Name:        "mapping",
			Destination: &c.Config.Mapping,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Destination: &c.Config.Debug,
		},
		cli.StringFlag{
			Name:        "config, c",
			Destination: &c.Config.ConfigFile,
		},
	}
}

func (c *CLI) defaultDmlFlags() []cli.Flag {
	return append(
		c.defaultFlags(),
		cli.StringFlag{
			Name:        "file, f",
			Destination: &c.Config.InputFile,
		},
		cli.StringFlag{
			Name:        "type, t",
			Destination: &c.Config.Type,
		},
		cli.StringFlag{
			Name:        "success-file",
			Value:       "./success.csv",
			Destination: &c.Config.SuccessPath,
		},
		cli.StringFlag{
			Name:        "error-file",
			Value:       "./error.csv",
			Destination: &c.Config.ErrorPath,
		},
	)
}
