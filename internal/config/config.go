package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	DBUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user"`
}

func (c *Config) SetUser(username string) (int, error) {
	c.CurrentUser = username

	err := Write(c)
	if err != nil {
		return fmt.Printf("Error: %v", err)
	}
	return 0, nil
}

func Read() (Config, error) {
	var config Config

	filePath, err := GetConfigPath()
	if err != nil {
		return config, fmt.Errorf("Error: %v", err)
	}

	configData, err := os.Open(filePath)
	defer configData.Close()
	if err != nil {
		return config, fmt.Errorf("Error: %v", err)
	}

	byteData, err := io.ReadAll(configData)
	if err != nil {
		return config, fmt.Errorf("Error: %v", err)
	}

	json.Unmarshal(byteData, &config)

	return config, nil

}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error: %v", err)
	}
	return homeDir + "/.gatorconfig.json", nil
}

func Write(conf *Config) error {
	data, err := json.Marshal(conf)

	filePath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	err = os.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	return nil
}
