package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// AppConfig from config yml
type AppConfig struct {
	Telegram *Telegram `yaml:"telegram"`
}

// Telegram is used to configure the Telegram service
type Telegram struct {
	HistoryLimit int `yaml:"history_limit"`
}

// Load config from file
func Load(filename string) (res *AppConfig, err error) {
	res = &AppConfig{}
	data, err := os.ReadFile(filename) // #nosec G304
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, res); err != nil {
		return nil, err
	}

	return res, nil
}
