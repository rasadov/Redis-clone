package main

import "strings"

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

func ReadConfig(filename string, storage *InMemoryStorage) error {
	data, err := readMapFromFile(filename)
	if err != nil {
		return err
	}
	storage.data = data
	return nil
}

func WriteConfig(filename string, data map[string]entry) error {
	return writeMapToFile(filename, data)
}
