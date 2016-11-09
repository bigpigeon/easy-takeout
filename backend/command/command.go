package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

type Config struct {
	BaseUrl      string            `"baseurl"`
	DbAddress    string            `"dbaddress"`
	DbType       string            `"dbtype"`
	NeedLogin    bool              `"needlogin"`
	DbArgs       map[string]string `"dbargs"`
	CacheAddress string            `"cacheaddress"`
}

type Command struct {
	Description string
	Executer    func(c *Config)
}

var AvaliableCommand = map[string]Command{}

var HelpDescription = `Usage of easy-takeout:
Usage:
  easy-takeout [flags]
  easy-takeout [commands] [flags]

`

func PrintCommands() {
	interval := 3

	maxNameLength := 0
	for name, _ := range AvaliableCommand {
		if len(name) > maxNameLength {
			maxNameLength = len(name)
		}
	}
	for name, c := range AvaliableCommand {
		spaces := strings.Repeat(" ", maxNameLength-len(name)+interval)
		fmt.Fprintln(os.Stderr, "  ", name, spaces, c.Description)
	}
}

func UsageGenerate() {
	var Usage = func() {
		fmt.Fprintf(os.Stderr, HelpDescription)
		fmt.Fprintf(os.Stderr, flag.Arg(0))
		fmt.Fprintf(os.Stderr, "Available Commands:\n")
		PrintCommands()
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Usage = Usage
}

//config get priority command > configFile
func Execute() {
	UsageGenerate()
	var commandExec func(c *Config)
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-") == false {
		ac := os.Args[1]
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
		if c, ok := AvaliableCommand[ac]; ok {
			commandExec = c.Executer
		}

	}
	if commandExec == nil {
		commandExec = AvaliableCommand["server"].Executer
	}

	var configFile string
	flag.StringVar(&configFile, "config", "", "config file(support type toml|yaml|json)")
	var baseUrl string
	flag.StringVar(&baseUrl, "baseurl", "", "base server address")
	var dbAddress string
	flag.StringVar(&dbAddress, "dbaddress", "", "database server address")
	var dbType string
	flag.StringVar(&dbType, "dbtype", "", "select mysql/sqlite3/postgresql")
	var cacheAddress string
	flag.StringVar(&cacheAddress, "cacheaddress", "", "cache server address")
	var needLogin bool
	flag.BoolVar(&needLogin, "needlogin", false, "only need name to takeout when needlogin was false")
	flag.Parse()

	config := &Config{}
	if configFile != "" {
		f, err := ioutil.ReadFile(configFile)
		if err == nil {
			if strings.HasSuffix(configFile, ".toml") {

				err = toml.Unmarshal(f, config)
				if err != nil {
					panic(err)
				}
			} else if strings.HasSuffix(configFile, ".json") {
				err = json.Unmarshal(f, config)
				if err != nil {
					panic(err)
				}
			} else if strings.HasSuffix(configFile, ".yaml") {
				err = yaml.Unmarshal(f, config)
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Println("目前不支持这种格式的配置文件")
			}
		} else {
			fmt.Errorf("%e, %s", err, "配置文件读取失败")
		}
	}
	if baseUrl != "" {
		config.BaseUrl = baseUrl
	}
	if dbAddress != "" {
		config.DbAddress = dbAddress
	}
	if dbType != "" {
		config.DbType = dbType
	}
	if needLogin == true {
		config.NeedLogin = needLogin
	}
	commandExec(config)
}
