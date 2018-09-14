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
	Username    string
	Password    string
	Endpoint    string
	ApiVersion  string
	Query       string
	Type        string
	InputFile   string
	Delimiter   string
	Encoding    string
	UpsertKey   string
	Output      string
	Format      string
	Mapping     string
	ErrorPath   string
	SuccessPath string
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
		Config: &config{},
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
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor()
				return executor.query(c.Config)
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
				executor := NewCommandExecutor()
				return executor.insert(c.Config)
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
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor()
				return executor.update(c.Config)
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
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor()
				return executor.upsert(c.Config)
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
					cli.StringFlag{
						Name:        "type, t",
						Destination: &c.Config.Type,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor()
				return executor.delete(c.Config)
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
					cli.StringFlag{
						Name:        "type, t",
						Destination: &c.Config.Type,
					},
				},
			),
			Action: func(ctx *cli.Context) error {
				executor := NewCommandExecutor()
				return executor.undelete(c.Config)
			},
		},
	}
	return app.Run(args)
}

func createCommandFlag(defaultFlags []cli.Flag, flags []cli.Flag) []cli.Flag {
	return append(defaultFlags, flags...)
}
