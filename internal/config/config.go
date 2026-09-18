package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return fmt.Sprintf("error with user home directory: %w", err), err
	}

	filePath := filepath.Join(homeDir, configFileName)

	return filePath, nil
}

func write(cfg *Config) error {
	path, err := getConfigFilePath()

	if err != nil {
		return err
	}

	marshalledJson, err := json.Marshal(cfg)

	err = os.WriteFile(path, marshalledJson, 0600)

	return err
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()

	if err != nil {
		return Config{}, fmt.Errorf("error opening config file: %w", err)
	}

	jsonData, err := os.ReadFile(filePath)

	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(jsonData, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func SetUser(userName *string, cfg *Config) error {

	if cfg.UserName == "" {
		currentUser, err := user.Current()

		if err != nil {
			return err
		}

		cfg.UserName = currentUser.Username
	}

	return write(cfg)
}
