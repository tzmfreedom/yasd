package cli

import (
	"fmt"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/urfave/cli"
)

type CLI struct {
	Config *config
	logger Logger
	client *soapforce.Client
}

type config struct {
	Username   string
	Password   string
	Endpoint   string
	ApiVersion string
	Query      string
	Type       string
	InputFile  string
	Delimiter  string
	Encoding   string
}

type Logger struct {
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
		client: soapforce.NewClient(),
	}
	return c
}

func (c *CLI) Run(args []string) error {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s revision=%s\n", c.App.Version, Revision)
	}

	defaultFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "username, u",
			Destination: &c.Config.Username,
			EnvVar:      "SF_USERNAME",
		},
		cli.StringFlag{
			Name:        "password, p",
			Destination: &c.Config.Password,
			EnvVar:      "SF_PASSWORD",
		},
		cli.StringFlag{
			Name:        "endpoint, e",
			Value:       "login.salesforce.com",
			Destination: &c.Config.Endpoint,
			EnvVar:      "SF_ENDPOINT",
		},
		cli.StringFlag{
			Name:        "apiversion",
			Value:       DefaultApiVersion,
			Destination: &c.Config.ApiVersion,
			EnvVar:      "SF_APIVERSION",
		},
		cli.StringFlag{
			Name:        "delimiter",
			Destination: &c.Config.Delimiter,
		},
		cli.StringFlag{
			Name:        "encoding",
			Destination: &c.Config.Encoding,
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
						Name:        "query",
						Destination: &c.Config.Query,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor()
				err := executor.query(c.Config)
				if err != nil {
					return err
				}
				return nil
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
						Name:        "file",
						Destination: &c.Config.InputFile,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor()
				err := executor.insert(c.Config)
				if err != nil {
					return err
				}
				return nil
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
						Name:        "file",
						Destination: &c.Config.InputFile,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor()
				err := executor.update(c.Config)
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "upsert",
			Usage: "Upsert SObject Record",
			Flags: createCommandFlag(
				defaultFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:        "file",
						Destination: &c.Config.InputFile,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor()
				err := executor.upsert(c.Config)
				if err != nil {
					return err
				}
				return nil
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
						Name:        "file",
						Destination: &c.Config.InputFile,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor()
				err := executor.delete(c.Config)
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "undelete",
			Usage: "Undelete SObject Record",
			Flags: createCommandFlag(
				defaultFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:        "file",
						Destination: &c.Config.InputFile,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor()
				err := executor.undelete(c.Config)
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
	app.Run(args)
	return nil
}

func createCommandFlag(defaultFlags []cli.Flag, flags []cli.Flag) []cli.Flag {
	return append(defaultFlags, flags...)
}
