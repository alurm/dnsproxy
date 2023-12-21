package proxy

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"path/filepath"
	"fmt"
)

type configurationJSON struct {
	Server    string
	Blocklist []string
}

type configuration struct {
	Server    net.IP
	Blocklist map[string]bool
}

var Conf configuration

func Here() string {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		panic(err)
	}
	return filepath.Dir(executable)
}

func ReadConfiguration(path string) configuration {
	var bytes []byte
	{
		file, err := os.Open(path)
		if err != nil {
			fmt.Println("configuration.json not found")
			os.Exit(1)
		}
		defer file.Close()
		bytes, err = io.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}

	var confJSON configurationJSON
	err := json.Unmarshal(
		bytes,
		&confJSON,
	)
	if err != nil {
		panic(err)
	}

	var conf configuration

	server := net.ParseIP(confJSON.Server)
	if server == nil {
		panic("server configuration is incorrect")
	}

	conf.Server = server

	conf.Blocklist = map[string]bool{}
	for _, domain := range confJSON.Blocklist {
		conf.Blocklist[domain+"."] = true
	}

	return conf
}
