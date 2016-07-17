package config

import (
	"fmt"
	"math"
)

func ExampleConfig_GetFloat64Ok() {
	c := New()

	c.Put("float32", float32(1.5))
	c.Put("float64", float64(2.5))
	c.Put("int", int(2))

	fmt.Println(c.GetFloat64Ok("float32"))
	fmt.Println(c.GetFloat64Ok("float64"))
	fmt.Println(c.GetFloat64Ok("int"))
	fmt.Println(c.GetFloat64Ok("does not exist"))
	//Output:
	//1.5 true
	//2.5 true
	//0 false
	//0 false
}

func ExampleConfig_GetInt64Ok() {
	c := New()

	c.Put("byte", byte(1))
	c.Put("uint8", uint8(2))
	c.Put("int8", int8(3))
	c.Put("uint16", uint16(4))
	c.Put("int16", int16(5))
	c.Put("uint32", uint32(6))
	c.Put("int32", int32(7))
	c.Put("int", int(8))
	c.Put("uint64", uint64(9))
	c.Put("int64", int64(10))
	c.Put("ooh", uint64(math.MaxUint64))
	c.Put("string", "foobar")

	fmt.Println(c.GetInt64Ok("byte"))
	fmt.Println(c.GetInt64Ok("uint8"))
	fmt.Println(c.GetInt64Ok("int8"))
	fmt.Println(c.GetInt64Ok("uint16"))
	fmt.Println(c.GetInt64Ok("int16"))
	fmt.Println(c.GetInt64Ok("uint32"))
	fmt.Println(c.GetInt64Ok("int32"))
	fmt.Println(c.GetInt64Ok("int"))
	fmt.Println(c.GetInt64Ok("uint64"))
	fmt.Println(c.GetInt64Ok("int64"))
	fmt.Println(c.GetInt64Ok("ooh"))
	fmt.Println(c.GetInt64Ok("string"))
	fmt.Println(c.GetInt64Ok("does not exist"))
	//Output:
	//1 true
	//2 true
	//3 true
	//4 true
	//5 true
	//6 true
	//7 true
	//8 true
	//9 true
	//10 true
	//-1 true
	//0 false
	//0 false
}

func ExampleConfig_GetValuesOk() {
	c := New()

	c.Put("a.b", 1)
	c.Put("a.c", 2)
	c.Put("d", 3)

	printlnOk := func(key string) {
		_, ok := c.GetValuesOk(key)
		fmt.Println(ok)
	}

	printlnOk("a")
	printlnOk("a.b")
	printlnOk("a.c")
	printlnOk("d")
	printlnOk("does not exist")
	//Output:
	//true
	//false
	//false
	//false
	//false
}
