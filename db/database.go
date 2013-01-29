package db

import (
	"io/ioutil"
)

type List interface{
	Dumps() string
	Loads(s string)
}

func Save(list List, filePath string) {
	bytes := list.Dumps()
	writeFile([]byte(bytes), filePath)
}

func Load(list List, filePath string) {
	bytes := readFile(filePath)
	list.Loads(string(bytes))
}

func writeFile(bytes []byte, filePath string) {
	err := ioutil.WriteFile(filePath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func readFile(filePath string) []byte{
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return bytes
}
