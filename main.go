package main

import (
	"strings"
	"os/exec"
	"bufio"
	"runtime"
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

func (c cmd) run(t string) {
	var cmd *exec.Cmd
	cm, flag := c.format(t)
	cmd = exec.Command(cm, flag)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}


func (c cmd) format(t string) (string, string) {
	t = strings.Trim(t, "\n")
	if runtime.GOOS == "windows" {
		return "cmd", fmt.Sprintf("/c %s %s.md", c, t)
	}
	return "sh", fmt.Sprintf("-s %s %s.md", c, t)
}

var conf = &Config{}
type Config struct {
	Path   string `toml:"path"`
	OutDir string `toml:"outdir"`
	Edit   cmd    `toml:"edit"`
	ServerPort int`toml:"serverport"`
	TemplateFile string `toml:"templatefile`
	assetsDir string `toml:"assetsdir"`
}

func (c *Config) read(p string){
	f, _ := os.Open(p)
	b, _ := ioutil.ReadAll(f)
	if _, err := toml.Decode(string(b), &conf); err != nil {
		panic(err)
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
		Action:  new,
	},
	{
		Name: "list",
		Aliases: []string{"l"},
		Usage:   "show note file list",
		Action:  list,
	},
	{
		Name: "remove",
		Aliases: []string{"r"},
		Usage:   "remove note file",
		Action:  remove,
	},
	{
		Name: "todo",
		Aliases: []string{"t"},
		Usage:   "todo list",
		Action:  todo,
	},
	{
		Name: "grep",
		Aliases: []string{"g"},
		Usage:   "grep file",
		Action:  grep,
	},
	{
		Name: "config",
		Aliases: []string{"c"},
		Usage:   "edit config file.",
		Action:  editConf,
	},
}

func new(c *cli.Context) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter Title:")
	t, _ := reader.ReadString('\n')
	conf.Edit.run(t)
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
func editConf(c *cli.Context) error {
	return nil
}

func main() {
	app := cli.NewApp()
	app.Commands = command
	app.Version = version
	app.Name = name
	conf.init()
	app.Run(os.Args)
}