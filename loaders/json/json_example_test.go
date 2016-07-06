package json_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/gogolfing/config"
	jsonloader "github.com/gogolfing/config/loaders/json"
)

func Example() {
	input := `{
		"string": "foo bar",
		"bool": false,
		"int": 12345678,
		"float": 1234.5678,
		"nested1": {
			"nested2": {
				"nested3": "..."
			}
		},
		"nope": null
	}`

	inputReader := strings.NewReader(input)

	loader := config.NewReaderFuncLoader(inputReader, (&jsonloader.Loader{}).LoadReader)

	c := config.New()
	_, err := c.PutLoaders(loader)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.Values().EachKeyValue(func(key config.Key, value interface{}) {
		log.Println(key, value)
	})

	fmt.Println(c.GetStringOk("string"))
	fmt.Println(c.GetBoolOk("bool"))
	fmt.Println(c.GetInt64Ok("int"))
	fmt.Println(c.GetFloat64Ok("float"))
	fmt.Println(c.GetStringOk("nested1.nested2.nested3"))
	fmt.Println(c.GetOk("nope"))

	fmt.Println()
	fmt.Println(c.GetInt64Ok("bool"))
	fmt.Println(c.GetOk("not present"))

	//Output:
	//foo bar true
	//false true
	//12345678 true
	//1234.5678 true
	//... true
	//<nil> true
	//
	//0 false
	//<nil> false
}
