package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

var test map[string]interface{}
var configFile []byte
var index int

var finalMap map[string]string

func init() {
	finalMap = map[string]string{}
	test = map[string]interface{}{}
}

func New(filePath string) {
	var err error

	configFile, err = ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("There was an error opening the config file!")
	}
}

func GetString(key, backup string) string {
	keys := parseKey(key)
	val := get(keys)
	if val == nil {
		return backup
	}
	return val.(string)
}

func get(keys []string) interface{} {
	var ret interface{}

	for _, k := range keys {
		if ret != nil {
			ret = ret.(map[string]interface{})[k]
			continue
		}
		ret = test[k]
	}
	return ret
}

func parseKey(key string) []string {
	keys := []string{}
	runes := []rune(key)

	buff := []rune{}
	for _, r := range runes {
		if r == '.' {
			keys = append(keys, string(buff))
			buff = []rune{}
			continue
		}
		buff = append(buff, r)
	}
	keys = append(keys, string(buff))

	return keys
}

func Parse() {
	json.Unmarshal(configFile, &test)
}

func Print() {
	fmt.Println(test)
}
