package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const (
	configFileName = ".gatorconfig.json"
	perm           = os.FileMode(0o755) // chmod rwxrw-rw-
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
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

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("unable to retrieve home directory")
	}
	return fmt.Sprintf("%s/%s", home, configFileName), nil
}

func write(cfg Config) error {
	cfgLocation, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.New("unable to marshal config")
	}
	if err := os.WriteFile(cfgLocation, data, perm); err != nil {
		return errors.New("unable to write config")
	}
	return nil
}
