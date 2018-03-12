package main

import (
	"time"
	"strings"
	"os/exec"
	"bufio"
	"runtime"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"fmt"
	"os"
	"text/template"
	"github.com/urfave/cli"
	"github.com/BurntSushi/toml"
)
const (
	version = "0.1"
	name = "whs"
	serverPort = 10290
)
var textTml = `# {{.Title}} 
`
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
	if runtime.GOOS == "windows" {
		return "cmd", fmt.Sprintf("/c %s %s", c, t)
	}
	return "sh", fmt.Sprintf("-s %s %s", c, t)
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
	p := filepath.Join(usr.HomeDir, ".whs.toml")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		c.createConfigFile(usr.HomeDir, p)
	}
	c.read(p)
}

func (c *Config) createConfigFile(h, p string){
	f, _ := os.Create(p)
	c.Path = p
	c.OutDir = filepath.Join(h, "whs", "_post")
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
	var f *os.File
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter Title:")

	os.MkdirAll(conf.OutDir, os.ModePerm)

	title, _ := reader.ReadString('\n')
	t := time.Now()
	title = strings.TrimSpace(title)
	fileName := filepath.Join(conf.OutDir, t.Format("2006-01-02_") + title + ".md")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		f, _ = os.Create(fileName)
	} else {
		fmt.Println("file is exist")
		os.Exit(0)
	}

	tml, _ := template.New("newPost").Parse(textTml)
	tml.Execute(f, struct{
		Title string
	}{
		Title: title,
	})

	conf.Edit.run(fileName)
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