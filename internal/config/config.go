package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fullPath := homeDir + "/" + configFileName
	return fullPath, nil

}

func write(cfg Config) error {

	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(fullPath, data, 0777)
	if err != nil {
		return err
	}

	return nil

}

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, err
	}

	config := Config{}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil

}

func (c *Config) SetUser(username string) {
	c.CurrentUsername = username
	err := write(*c)
	if err != nil {
		return
	}

}
