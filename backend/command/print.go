package command

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func printConfig(c *Config) {
	s := reflect.ValueOf(c).Elem()
	typeOfT := s.Type()
	NameValuePairs := map[string]interface{}{}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		NameValuePairs[typeOfT.Field(i).Name] = f.Interface()
	}
	b, err := json.MarshalIndent(NameValuePairs, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func init() {
	AvaliableCommand["print"] = Command{
		Description: "print final config",
		Executer:    printConfig,
	}
}
