package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const (
	ExitCodeOK int = iota
	ExitCodeError

	AppName           = "yasd"
	Usage             = "Yet Another Salesforce Dataloader"
	DefaultApiVersion = "38.0"
)

var (
	Version  string
	Revision string
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s revision=%s\n", c.App.Version, Revision)
	}

	app := cli.NewApp()
	app.Name = AppName
	app.Usage = Usage
	app.Version = Version
	app.Commands = Commands

	err := app.Run(os.Args)

	statusCode := ExitCodeOK
	if err != nil {
		fmt.Println(err.Error())
		statusCode = ExitCodeError
	}
	os.Exit(statusCode)
}
