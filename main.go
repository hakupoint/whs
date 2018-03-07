package main

import (
	"io/ioutil"
	"os/user"
	"path/filepath"
	"fmt"
	"os"
	"github.com/urfave/cli"
	"github.com/BurntSushi/toml"
)
const (
	version = "0.1"
	name = "whs"
	serverPort = 10290
)
type cmd string

var conf = &Config{}
type Config struct {
	Path   string
	OutDir string
	Edit   cmd
	ServerPort int
	TemplateFile string
}

func (c *Config) read(p string){
	f, _ := os.Open(p)
	b, _ := ioutil.ReadAll(f)
	if _, err := toml.Decode(string(b), &conf); err == nil {
		fmt.Printf("%+v", conf)
	}
}

func (c *Config) init() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	p := filepath.Join(usr.HomeDir, "whs.toml")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		c.createConfigFile(p)
	}
	c.read(p)
}

func (c *Config) createConfigFile(p string){
	f, _ := os.Create(p)
	c.Path = p
	c.OutDir = filepath.Join(p, "_post")
	c.ServerPort = serverPort
	c.Edit = "vim"
	e := toml.NewEncoder(f)
	err := e.Encode(c)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
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
		Name: "remove",
		Aliases: []string{"r"},
		Usage:   "remove note file",
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
	{
		Name: "config",
		Aliases: []string{"c"},
		Usage:   "edit config file.",
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
func remove(c *cli.Context) error {
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
	app.Commands = command
	app.Version = version
	app.Name = name
	app.Run(os.Args)
	conf.init()
}