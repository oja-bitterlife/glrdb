package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	DBName    string   `toml:"db_name"` // データベースファイル名
	Sources   []Source `toml:"sources"`
	Blacklist []string `toml:"blacklist"` // node_modules などをここに追加予定
	MaxDepth  int      `toml:"max_depth"` // デフォルトは64、必要に応じて変更可能
}

type Source struct {
	Path string `toml:"path"`
}

func newDefaultConfig() *Config {
	return &Config{
		DBName:   "glrdb.boltdb",
		MaxDepth: 64,
		// とりあえずデフォルトを手入力
		Blacklist: []string{"node_modules", "vendor"},
	}
}

func loadConfig(path string) (*Config, error) {
	config := newDefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
