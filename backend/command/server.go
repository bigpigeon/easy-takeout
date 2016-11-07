package command

import (
	"net/url"

	"github.com/easy-takeout/easy-takeout/backend/api"

	"github.com/easy-takeout/easy-takeout/backend/definition"
)

func server(c *Config) {
	db, err := definition.Connect(c.DbType, c.DbAddress, c.DbArgs)
	if err != nil {
		panic(err)
	}
	url, err := url.Parse(c.BaseUrl)
	if err != nil {
		panic(err)
	}
	route := api.Create(db, c.NeedLogin)
	route.Run(url.Host)
}

func init() {
	AvaliableCommand["server"] = Command{
		Description: "run webserver",
		Executer:    server,
	}
}
