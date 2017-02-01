package command

import (
	"net/url"

	"github.com/bigpigeon/easy-takeout/backend/api"

	"log"

	"time"

	"github.com/bigpigeon/easy-takeout/backend/cachemanage"
	"github.com/bigpigeon/easy-takeout/backend/definition"
	"github.com/fsnotify/fsnotify"
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
	generate(c)
	// run generate when template file modify
	if c.Watch == true {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		go func() {
			needReload := false
			t := time.NewTicker(100 * time.Millisecond)
			defer t.Stop()
			for {
				select {
				case <-watcher.Events:
					needReload = true
				case <-t.C:
					if needReload == true {
						log.Println("html reloading...")
						generate(c)
						needReload = false
					}
				case err := <-watcher.Errors:
					log.Println("error:", err)
					break
				}
			}
		}()
		watcher.Add("template")
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
