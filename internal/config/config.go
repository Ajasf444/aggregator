package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
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
	// TODO: use file, err := os.Open() and defer file.Close() and decoder := json.NewDecoder() and decoder.Decode(&cfg)
	cfgLocation, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	content, err := os.ReadFile(cfgLocation)
	if err != nil {
		return Config{}, errors.New("unable to read config file")
	}
	cfg := Config{}
	if err := json.Unmarshal(content, &cfg); err != nil {
		return Config{}, errors.New("unable to unmarshal config file")
	}
	return cfg, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("unable to retrieve home directory")
	}
	return filepath.Join(home, configFileName), nil
}

func write(cfg Config) error {
	// TODO: use file, err := os.Create() and defer file.Close() and encoder := json.NewEncoder() and encoder.Encode(cfg)
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
