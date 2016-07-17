package env

import (
	"fmt"
	"os"

	"github.com/gogolfing/config"
)

func ExampleNewPrefixLowerUnderscoreLoader() {
	prefix := "ANW8THWNFQ4IWH874H__" //does not matter at all what this is. just don't want conflicts while testing.

	os.Setenv(prefix+"ONE_Two", "value")
	os.Setenv("not-my-prefix", "do not want")

	loader := NewPrefixLowerUnderscoreLoader(prefix)

	c := config.New()

	_, err := c.MergeLoaders(loader)
	if err != nil {
		fmt.Println(err)
	}

	want := config.New()
	want.Put("one.two", "value")

	if !c.EqualValues(want) {
		fmt.Println("we didn't get what we wanted")
	}

	fmt.Println(c.GetStringOk("one.two"))
	//Output:
	//value true
}
