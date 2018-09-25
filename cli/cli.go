package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

type CLI struct {
	Config *config
}

type config struct {
	Username         string
	Password         string
	Endpoint         string
	ApiVersion       string
	Query            string
	Type             string
	InputFile        string
	Delimiter        string
	Encoding         string
	UpsertKey        string
	Output           string
	Format           string
	Mapping          string
	ErrorPath        string
	SuccessPath      string
	EncyptionKeyPath string
	BatchSize        int
	Debug            bool
	InsertNulls      bool
}

var (
	Version  string
	Revision string
)

const (
	AppName = "yasd"
	Usage   = ""
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
	defaultEncryptionKeyPath, err := defaultEncryptionKeyPath()
	if err != nil {
		return err
	}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s revision=%s\n", c.App.Version, Revision)
	}

	defaultFlags := []cli.Flag{
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
			Name:        "apiversion",
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
			Flags: createCommandFlag(
				defaultFlags,
				[]cli.Flag{
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
			Flags: createCommandFlag(
				defaultFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:        "file, f",
						Destination: &c.Config.InputFile,
					},
					cli.StringFlag{
						Name:        "type, t",
						Destination: &c.Config.Type,
					},
					cli.StringFlag{
						Name:        "successfile",
						Value:       "./success.csv",
						Destination: &c.Config.SuccessPath,
					},
					cli.StringFlag{
						Name:        "errorfile",
						Value:       "./error.csv",
						Destination: &c.Config.ErrorPath,
					},
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
			Flags: createCommandFlag(
				defaultFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:        "file, f",
						Destination: &c.Config.InputFile,
					},
					cli.StringFlag{
						Name:        "type, t",
						Destination: &c.Config.Type,
					},
					cli.BoolFlag {
						Name:        "insert-nulls",
						Destination: &c.Config.InsertNulls,
					},
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
			Flags: createCommandFlag(
				defaultFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:        "file, f",
						Destination: &c.Config.InputFile,
					},
					cli.StringFlag{
						Name:        "type, t",
						Destination: &c.Config.Type,
					},
					cli.StringFlag{
						Name:        "upsertkey, k",
						Value:       "Id",
						Destination: &c.Config.UpsertKey,
					},
					cli.BoolFlag {
						Name:        "insert-nulls",
						Destination: &c.Config.InsertNulls,
					},
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
			Flags: createCommandFlag(
				defaultFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:        "file, f",
						Destination: &c.Config.InputFile,
					},
					cli.StringFlag{
						Name:        "type, t",
						Destination: &c.Config.Type,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.delete(*c.Config)
			},
		},
		{
			Name:  "undelete",
			Usage: "Undelete SObject Record",
			Flags: createCommandFlag(
				defaultFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:        "file, f",
						Destination: &c.Config.InputFile,
					},
					cli.StringFlag{
						Name:        "type, t",
						Destination: &c.Config.Type,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.undelete(*c.Config)
			},
		},
		{
			Name:  "generate-key",
			Usage: "Generate AES Key",
			Flags: createCommandFlag(
				defaultFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:        "key",
						Value:       defaultEncryptionKeyPath,
						Destination: &c.Config.EncyptionKeyPath,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.generateEncryptionKey(*c.Config)
			},
		},
		{
			Name:  "decrypt-key",
			Usage: "Descrypt AES Key",
			Flags: createCommandFlag(
				defaultFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:        "key",
						Destination: &c.Config.EncyptionKeyPath,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor(c.Config.Debug)
				return executor.debug(*c.Config)
			},
		},
	}
	return app.Run(args)
}

func createCommandFlag(defaultFlags []cli.Flag, flags []cli.Flag) []cli.Flag {
	return append(defaultFlags, flags...)
}

func configDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "yasd"), nil
}

func defaultEncryptionKeyPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "key"), nil
}

func createEncryptionKeyFile(filename string) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err = os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	configPath := filepath.Join(dir, filename)
	if _, err := os.Stat(configPath); err != nil {
		if err = generateEncryptionKey(configPath); err != nil {
			return err
		}
	}
	return nil
}
