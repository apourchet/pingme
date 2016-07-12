package ping

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/user"
	"path"
)

type ClientConfig struct {
	Host    string            `json:"host"`
	Port    string            `json:"port"`
	Aliases map[string]string `json:"aliases"`
}

const (
	CLIENT_CONFIG       = ".pingme"
	DEFAULT_SERVER_PORT = "1025"

	DEFAULT_HOST = "localhost"
	DEFAULT_PORT = "1025"
)

func GetClientConfig() ClientConfig {
	conf := ClientConfig{DEFAULT_HOST, DEFAULT_PORT, make(map[string]string)}
	homeDir, err := getHomeDir()

	if err != nil {
		log.Println("Could not load config file. Using Defaults")
		return conf
	}

	fullPath := path.Join(homeDir, CLIENT_CONFIG)
	fileBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		conf.WriteOut()
		log.Println("Created config file at " + fullPath)
		return conf
	}

	err = json.Unmarshal(fileBytes, &conf)
	if err != nil {
		log.Println("Format error in config file. Erase it to use defaults.")
		return conf
	}
	return conf
}

func (conf *ClientConfig) WriteOut() error {
	homeDir, err := getHomeDir()
	if err != nil {
		return err
	}

	fileBytes, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return err
	}

	fullPath := path.Join(homeDir, CLIENT_CONFIG)
	err = ioutil.WriteFile(fullPath, fileBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
