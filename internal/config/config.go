package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	// set the username in the struct
	if c.CurrentUserName == username {
		return nil
	}
	c.CurrentUserName = username

	// write back to file
	filePath, err := GetConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&c); err != nil {
		return err
	}

	return nil
}

func GetConfigPath() (string, error) {
	const configPath = ".gatorconfig.json"
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(homePath, configPath)
	return filePath, nil
}

func Read() (*Config, error) {

	var config *Config

	filePath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
