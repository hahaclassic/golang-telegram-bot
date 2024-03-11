package telegram

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MainToken   string `yaml:"main-tg-bot-token"`
	LoggerToken string `yaml:"logger-tg-bot-token"`
	AdminChatID int    `yaml:"admin"`
	StoragePath string `yaml:"storage-path"`
	Host        string `yaml:"host"`
	BatchSize   int    `yaml:"batch-size"`
	ErrChanSize int    `yaml:"err-chan-size"`
}

func NewConfig(configPath string) (*Config, error) {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
