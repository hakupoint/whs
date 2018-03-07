package main

import (
	"os"
	"github.com/urfave/cli"
)
type cmd string

type Config struct {
	Path   string
	OutDir cmd
	Edit   cmd
}

var command = []cli.Command{
	{
		Name: "new",
		Aliases: []string{"n"},
		Usage:   "add new note file",
		Action:  func(c *cli.Context) error {
			return nil
		},
	},
	{
		Name: "list",
		Aliases: []string{"l"},
		Usage:   "show note file list",
		Action:  func(c *cli.Context) error {
			return nil
		},
	},
	{
		Name: "todo",
		Aliases: []string{"t"},
		Usage:   "todo list",
		Action:  func(c *cli.Context) error {
			return nil
		},
	},
	{
		Name: "grep",
		Aliases: []string{"g"},
		Usage:   "grep file",
		Action:  func(c *cli.Context) error {
			return nil
		},
	},
}

func add(c *cli.Context) error {
	return nil
}
func list(c *cli.Context) error {
	return nil
}
func todo(c *cli.Context) error {
	return nil
}
func grep(c *cli.Context) error {
	return nil
}

func main() {
	app := cli.NewApp()
	app.Run(os.Args)
}