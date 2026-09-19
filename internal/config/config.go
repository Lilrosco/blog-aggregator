package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(userName string) error {
	if userName == "" && cfg.CurrentUserName == "" {
		currentUser, err := user.Current()

		if err != nil {
			return err
		}

		cfg.CurrentUserName = currentUser.Username
	}

	cfg.CurrentUserName = userName

	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	filePath := filepath.Join(homeDir, configFileName)
	return filePath, nil
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()

	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)

	if err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()

	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(fullPath)
	
	if err != nil {
		return Config{}, err
	}

	defer file.Close()

	cfg := Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)

	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
