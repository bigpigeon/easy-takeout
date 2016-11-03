package command

import (
	"github.com/easy-takeout/easy-takeout/backend/api"
)

func server(c *Config) {

}

func init() {
	AvaliableCommand["server"] = Command{
		Description: "run webserver",
		Executer:    server,
	}
}
