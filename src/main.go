package main

import (
	"fmt"
)

// バージョンは頻繁に書き換えるので個別で定数化しておく
const version = "0.1.0"

const (
	// コンフィグのファイル名
	configFileName = "glrdb.toml"
)

type Repository struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func main() {
	config, err := loadConfig(configFileName)
	if err != nil {
		panic(err)
	}

	allRepos, err := scanDir(config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("--- start fetching description ---\n")

	if err = updateDB(config, allRepos); err != nil {
		panic(err)
	}
}
