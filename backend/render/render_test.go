package render

import (
	"testing"
)

func TestGetHtmlFile(t *testing.T) {
	source := "test/template"
	exclude := []string{"css"}
	files, keeps, err := SplitRenderFile(source, exclude)
	if err != nil {
		t.Error(err)
	}
	for _, f := range files {
		if f.Path == "css" {
			t.Fail()
		}
	}
	for _, f := range keeps {
		if f.Path == "example" || f.Path == "js" {
			t.Fail()
		}
	}
}

func TestRender(t *testing.T) {
	source := "test/template"
	target := "test/public"
	exclude := []string{"css", "js"}
	data := map[string]string{
		"name": "pigeon",
	}
	err := Render(source, target, exclude, data)
	if err != nil {
		t.Error(err)
	}
}
