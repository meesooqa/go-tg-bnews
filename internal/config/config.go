package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// AppConfig from config yml
type AppConfig struct {
	Log *Log `yaml:"log"`
}

// Log - log parameters
type Log struct {
	LevelCode    string `yaml:"level_code"`
	Level        slog.Level
	OutputFormat string `yaml:"output_format"`
	Path         string `yaml:"path"`
}

// load config from file
func load(fname string) (res *AppConfig, err error) {
	res = &AppConfig{}
	data, err := os.ReadFile(fname) // #nosec G304
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(data, res); err != nil {
		return nil, err
	}

	level, ok := logLevelMap[res.Log.LevelCode]
	if ok {
		res.Log.Level = level
	} else {
		res.Log.Level = slog.LevelInfo
	}

	return res, nil
}
