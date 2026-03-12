package main

import (
	// bolt "go.etcd.io/bbolt"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2" // /v2 を忘れずに
)

// バージョンは頻繁に書き換えるので個別で定数化しておく
const version = "0.1.0"

const (
	// データベースのファイル名
	dbFileName = "glrdb.db"
)

type Config struct {
	Paths []string `toml:"path"`
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(dbFileName)
	if err != nil {
		return nil, err
	}

	var config Config
	// BurntSushi版とほぼ同じシグネチャで使えます
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Loaded paths: %v\n", config.Paths)
}
