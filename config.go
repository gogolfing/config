package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

var data map[string]interface{}
var configFile []byte

func init() {
	data = map[string]interface{}{}
}

func New(filePath string) {
	var err error

	configFile, err = ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("There was an error opening the config file!")
	}
}

func Parse() {
	json.Unmarshal(configFile, &data)
}

func GetString(key, backup string) string {
	val := get(key)

	if val == nil {
		return backup
	}

	return val.(string)
}

func GetInt(key string, backup int) int {
	val := get(key)

	if val == nil {
		return backup
	}

	return int(val.(float64))
}

func GetBool(key string, backup bool) bool {
	val := get(key)

	if val == nil {
		return backup
	}

	return val.(bool)
}

func GetFloat64(key string, backup float64) float64 {
	val := get(key)

	if val == nil {
		return backup
	}

	return val.(float64)
}

func get(key string) interface{} {
	keys := strings.Split(key, ".")
	var ret interface{}

	for _, k := range keys {
		if ret != nil {
			ret = ret.(map[string]interface{})[k]
			continue
		}
		ret = data[k]
	}

	return ret
}
