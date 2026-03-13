package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Sources   []Source `toml:"sources"`
	Blacklist []string `toml:"blacklist"` // node_modules などをここに追加予定
}

type Source struct {
	Path string `toml:"path"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// とりあえず手入力しておく
	if len(config.Blacklist) == 0 {
		config.Blacklist = []string{"node_modules", "vendor"}
	}

	return &config, nil
}
