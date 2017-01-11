package command

import (
	"net/url"

	"github.com/bigpigeon/easy-takeout/backend/api"

	"github.com/bigpigeon/easy-takeout/backend/cachemanage"
	"github.com/bigpigeon/easy-takeout/backend/definition"
)

func server(c *Config) {
	db, err := definition.Connect(c.DbType, c.DbAddress, c.DbArgs)
	if err != nil {
		panic(err)
	}
	var cacheClient *cachemanage.Manage
	if c.CacheAddress == "" {
		cacheClient, err = cachemanage.Create(c.DbType, c.DbAddress, c.DbArgs)
	} else {
		cacheClient, err = cachemanage.Create("redis", c.CacheAddress, nil)
	}
	if err != nil {
		panic(err)
	}

	url, err := url.Parse(c.BaseUrl)
	if err != nil {
		panic(err)
	}
	route := api.Create(db, cacheClient, c.NeedLogin)
	route.Run(url.Host)
}

func init() {
	AvaliableCommand["server"] = Command{
		Description: "run webserver",
		Executer:    server,
	}
}
