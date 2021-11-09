package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var Config *Configuration

type Configuration struct {
	// Server配置
	Name    string
	Host    string
	TcpPort int

	// Kinx配置
	Version    string
	MaxConn    int
	MaxPackage uint32
}

func isExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func (c *Configuration) reload() {
	filePath := "conf/kinx.json"
	if isExist(filePath) {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("read config.json fail")
		}
		err = json.Unmarshal(data, &Config)
		if err != nil {
			fmt.Println("parse config.json fail")
		}
	}
}

func init() {
	Config = &Configuration{
		Name:       "TcpServerApp",
		Host:       "0.0.0.0",
		TcpPort:    8888,
		Version:    "0.4",
		MaxConn:    3,
		MaxPackage: 1024,
	}
	Config.reload()

	fmt.Printf("load configuration: %#v\n", Config)
}
