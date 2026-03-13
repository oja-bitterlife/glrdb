package main

import (
	"fmt"

	"go.etcd.io/bbolt"
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

	db, err := bbolt.Open(config.DBName, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	allRepos, err := scanDir(config)
	if err != nil {
		panic(err)
	}

	// debug
	fmt.Printf("Total repositories found: %#v\n", allRepos)

	for _, repoPath := range allRepos {
		fetchReadme(repoPath)
	}
}
