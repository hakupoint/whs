package main

import (
	"sync"
	"sort"
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

var mu sync.WaitGroup
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

type FileList []os.FileInfo

func (f FileList) Len() int {
	return len(f)
}

func (f FileList) Less(i, j int) bool {
	b := f[i].ModTime().Sub(f[j].ModTime())
	if b > 0 {
		return true
	} else {
		return false
	}
}

func (f FileList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type Results struct {
	Name string
	LineNo int
	LineContext string
}

func (r Results) string() string{
	return fmt.Sprintf("%d %s %s\n", r.LineNo, r.Name, r.LineContext)
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
	var fs FileList
	fs, _ = ioutil.ReadDir(conf.OutDir)
	sort.Sort(fs)
	for _, f := range fs {
		fmt.Println(f.Name())
	}
	return nil
}
func remove(c *cli.Context) error {
	return nil
}
func todo(c *cli.Context) error {
	return nil
}
func grep(c *cli.Context) error {
	var word = c.Args().Get(0)
	fs, _ := ioutil.ReadDir(conf.OutDir)
	for _, f := range fs {
		go func(f os.FileInfo) {
			fi, _ := os.Open(filepath.Join(conf.OutDir, f.Name()))
			scan := bufio.NewScanner(fi)
			index := 1
			for scan.Scan() {
				b := strings.Index(scan.Text(), word)
				if b != -1 {
					fmt.Print(Results{
						Name: f.Name(),
						LineContext: scan.Text(),
						LineNo: index,
					}.string())
				}
				index++
			}
			mu.Done()
		}(f)
		mu.Add(1)
	}
	mu.Wait();
	return nil
}
func editConf(c *cli.Context) error {
	conf.Edit.run(conf.Path)
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