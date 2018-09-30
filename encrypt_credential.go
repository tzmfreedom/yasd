package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

func encryptCredential(c *cli.Context) error {
	if err := validateEncryptCredential(c); err != nil {
		return err
	}
	fmt.Print("Password: ")
	b, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return err
	}
	fmt.Println("")
	password := string(b)
	if password == "" {
		return cli.NewExitError("password is not blank", 1)
	}
	keypath := c.String("key")
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

func validateEncryptCredential(c *cli.Context) error {
	key := c.String("key")
	if key == "" {
		_ = cli.ShowCommandHelp(c, "encrypt")
		return cli.NewExitError("key is required", 1)
	}
	if _, err := os.Stat(key); err != nil {
		_ = cli.ShowCommandHelp(c, "encrypt")
		return cli.NewExitError(fmt.Sprintf("No such file or directory: %s", key), 1)
	}
	return nil
}
