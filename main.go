package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

var mu sync.WaitGroup

const (
	version    = "0.1"
	name       = "whs"
	serverPort = 10290
)

var textTml = `# {{.Title}} 
`

type cmd string

func (c cmd) run(t string) {
	var cmd *exec.Cmd
	cm, param, flag := c.format(t)
	fmt.Println(cm, flag)
	cmd = exec.Command(cm, param, flag)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("sdfdsf")
		fmt.Println(err)
	}
}

func (c cmd) format(t string) (string, string, string) {
	if runtime.GOOS == "windows" {
		return "cmd", "/c ", fmt.Sprintf("%s %s", c, t)
	}
	return "sh", "-c", fmt.Sprintf("%s %s", c, t)
}

var conf = &Config{}

type Config struct {
	Path         string `toml:"path"`
	OutDir       string `toml:"outdir"`
	Edit         cmd    `toml:"edit"`
	ServerPort   int    `toml:"serverport"`
	TemplateFile string `toml:"templatefile`
	assetsDir    string `toml:"assetsdir"`
}

func (c *Config) read(p string) {
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

func (c *Config) createConfigFile(h, p string) {
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
	}
	return false
}

func (f FileList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type Results struct {
	Name        string
	Line        []struct{
		LineNo int
		LineContext string
	}
}

func (r Results) Print() {
	if len(r.Line) > 0 {
		fmt.Printf("\n%s\n", r.Name)
	}
	for _, v := range r.Line {
		fmt.Printf("%d:-> %s\n", v.LineNo, v.LineContext)
	}
}

var command = []cli.Command{
	{
		Name:    "new",
		Aliases: []string{"n"},
		Usage:   "add new note file",
		Action:  new,
	},
	{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "show note file list",
		Action:  list,
	},
	{
		Name:    "remove",
		Aliases: []string{"r"},
		Usage:   "remove note file",
		Action:  remove,
	},
	{
		Name:    "todo",
		Aliases: []string{"t"},
		Usage:   "todo list",
		Action:  todo,
	},
	{
		Name:    "grep",
		Aliases: []string{"g"},
		Usage:   "grep file",
		Action:  grep,
	},
	{
		Name:    "config",
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
	fileName := filepath.Join(conf.OutDir, t.Format("2006-01-02_")+title+".md")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		f, _ = os.Create(fileName)
	} else {
		conf.Edit.run(fileName)
		return nil
	}

	tml, _ := template.New("newPost").Parse(textTml)
	tml.Execute(f, struct {
		Title string
	}{
		Title: title,
	})
	defer f.Close()

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
	var results []Results
	if word == "" {
		os.Exit(0)
	}
	fs, _ := ioutil.ReadDir(conf.OutDir)
	for _, f := range fs {
		mu.Add(1)
		go func(f os.FileInfo) {
			fi, _ := os.Open(filepath.Join(conf.OutDir, f.Name()))
			scan := bufio.NewScanner(fi)
			index := 1
			var result Results
			
			for scan.Scan(){
				b := strings.Index(scan.Text(), word)
				if b != -1 {
					result.Name = f.Name()
					result.Line = append(result.Line, struct{
						LineNo int
						LineContext string
					}{
						LineNo: index,
						LineContext: scan.Text(),
					})
				}
				index++
			}
			results = append(results, result)
			mu.Done()
		}(f)
	}
	mu.Wait()
	for _, resu := range results {
		resu.Print()
	}
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
