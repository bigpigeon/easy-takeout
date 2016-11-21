package command

import (
	"github.com/easy-takeout/easy-takeout/backend/render"
)

func generate(c *Config) {
	err := render.Render("template", "public", []string{"static"}, c)
	if err != nil {
		panic(err)
	}
}

func init() {
	AvaliableCommand["generate"] = Command{
		Description: "generate static html with template",
		Executer:    generate,
	}
}
