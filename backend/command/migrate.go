package command

import (
	"github.com/bigpigeon/easy-takeout/backend/definition"
)

func migrate(c *Config) {
	db, err := definition.Connect(c.DbType, c.DbAddress, c.DbArgs)
	if err != nil {
		panic(err)
	}
	err = definition.Migrate(db)
	if err != nil {
		panic(err)
	}
}

func init() {
	AvaliableCommand["migrate"] = Command{
		Description: "migrate or create table",
		Executer:    migrate,
	}
}
