package flag

import (
	"fmt"

	"github.com/gogolfing/config"
)

func Example() {
	loader := New("")

	boolVar := false

	loader.FlagSet.String("s", "", "")
	loader.FlagSet.Int("i", 0, "")
	loader.FlagSet.BoolVar(&boolVar, "b", false, "")
	loader.FlagSet.Float64("f", 0, "")
	loader.FlagSet.String("withdefault", "default", "")
	loader.FlagSet.String("one-two-three", "", "")
	loader.FlagSet.String("aliased", "", "")

	loader.AddAlias("aliased", "a")

	//This means that withdefault -> default will be inserted in the resulting Config.
	//Leaving this blank would not insert this associated because the -withdefault
	//flag is not set below.
	loader.LoadDefaults = true

	loader.Args = []string{
		"-s", "sValue",
		"-i", "12",
		"-b=true",
		"-f", "1.2",
		//notice no withdefault flag set
		"-one-two-three", "three",
		"-aliased", "aValue",
	}

	c := config.New()

	_, err := c.MergeLoaders(loader)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(c.GetStringOk("s"))
	fmt.Println(c.GetInt64Ok("i"))
	fmt.Println(c.GetBoolOk("b"))
	fmt.Println(c.GetFloat64Ok("f"))
	fmt.Println(c.GetStringOk("withdefault"))
	fmt.Println(c.GetStringOk("one.two.three"))
	fmt.Println(c.GetStringOk("a")) //notice the use of the a key - aliased above.
	//Output:
	//sValue true
	//12 true
	//true true
	//1.2 true
	//default true
	//three true
	//aValue true
}
