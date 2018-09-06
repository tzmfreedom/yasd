package main

type CLI struct {
	Confif *config
	logger Logger
}

type config struct {
}

type Logger struct {
}

var (
	Version  string
	Revision string
)

const (
	DEFAULT_API_VERSION = "38.0"
)

func NewCli() *CLI {
	c := &CLI{}
	return c
}

func (*CLI) Run(args []string) error {
	return nil
}
