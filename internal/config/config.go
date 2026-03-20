package config

import (
	"encoding/json"
	"errors"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read(location string) (Config, error) {
	cfgLocation, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	content, err := os.ReadFile(cfgLocation)
	if err != nil {
		return Config{}, errors.New("unable to read config file")
	}
	config := Config{}
	if err := json.Unmarshal(content, &config); err != nil {
		return Config{}, errors.New("unable to unmarshal config file")
	}
	return config, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("unable to retrieve home directory")
	}
	return home + configFileName, nil
}

func SetUser() {
	// TODO:Marshal Config into JSON and then os.WriteFile()
}
