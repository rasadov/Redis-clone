package main

import (
	"strings"
)

var Config map[string]string

func InitConfig(args []string) {
	Config = make(map[string]string)
	isVal := false
	var key string

	for _, val := range args {
		if isVal {
			Config[key] = val
			isVal = false
		} else if strings.HasPrefix(val, "--") {
			key = strings.TrimPrefix(val, "--")
			isVal = true
		}
	}
}
