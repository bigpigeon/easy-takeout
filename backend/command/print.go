package command

import (
	"os"

	"github.com/BurntSushi/toml"
)

func printConfig(c *Config) {
	encoder := toml.NewEncoder(os.Stdout)
	err := encoder.Encode(c)
	if err != nil {
		panic(err)
	}
}

func init() {
	AvaliableCommand["print"] = Command{
		Description: "print final config",
		Executer:    printConfig,
	}
}
