package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
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
	keys := strings.Split(key, ".")
	val := get(keys)

	if val == nil {
		return backup
	}

	return val.(string)
}

func GetInt(key string, backup int) int {
	keys := strings.Split(key, ".")
	val := get(keys)

	if val == nil {
		return backup
	}

	return int(val.(float64))
}

func GetBool(key string, backup bool) bool {
	keys := strings.Split(key, ".")
	val := get(keys)

	if val == nil {
		return backup
	}

	return val.(bool)
}

func GetFloat64(key string, backup float64) float64 {
	keys := strings.Split(key, ".")
	val := get(keys)

	if val == nil {
		return backup
	}

	return val.(float64)
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

func Parse() {
	json.Unmarshal(configFile, &test)
}

func Print() {
	fmt.Println(test)
}
