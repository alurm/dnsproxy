package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type configuration struct {
	Server    string
	Blacklist []string
}

func readConfiguration() configuration {
	var c configuration
	var bytes []byte
	{
		executable, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		executable, err = filepath.EvalSymlinks(executable)
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Dir(executable) + "/configuration.json"
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		bytes, err = io.ReadAll(file)
	}

	err := json.Unmarshal(
		bytes,
		&c,
	)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func main() {
	c := readConfiguration()
	fmt.Printf("%#v\n", c)
}
