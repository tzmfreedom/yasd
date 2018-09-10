package main

import (
	"os"

	"github.com/tzmfreedom/yasd/cli"
)

const (
	ExitCodeOK int = iota
	ExitCodeError
)

func main() {
	cli := cli.NewCli()
	err := cli.Run(os.Args)

	statusCode := ExitCodeOK
	if err != nil {
		statusCode = ExitCodeError
	}
	os.Exit(statusCode)
}
