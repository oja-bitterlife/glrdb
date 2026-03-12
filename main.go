package main

import (
	// bolt "go.etcd.io/bbolt"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2" // /v2 を忘れずに
)

// バージョンは頻繁に書き換えるので個別で定数化しておく
const version = "0.1.0"

const (
	// データベースのファイル名
	dbFileName = "glrdb.db"
	// コンフィグのファイル名
	configFileName = "glrdb.toml"
)

type Config struct {
	Sources []Source `toml:"sources"`
}

type Source struct {
	Path string `toml:"path"`
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configFileName)
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

func isBareRepo(path string) bool {
	_, errObjects := os.Stat(filepath.Join(path, "objects"))
	_, errRefs := os.Stat(filepath.Join(path, "refs"))
	_, errConfig := os.Stat(filepath.Join(path, "config"))
	_, errHead := os.Stat(filepath.Join(path, "HEAD"))
	return errObjects == nil && errRefs == nil && errConfig == nil && errHead == nil
}
func scanRepositories(config *Config) error {
	for _, src := range config.Sources {
		fmt.Printf("Scanning: %s\n", src.Path)

		// ここで走査を開始
		filepath.WalkDir(src.Path, func(path string, d os.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}

			// 高速化のため、いくつかのディレクトリをスキップ
			// ----------------------------------------
			name := d.Name()

			// ドットで始まるディレクトリはスキップ
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			// node_modulesとか、普通のディレクトリ名でドンと居座るのやめて(´・ω:;.:...
			if name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}

			// Bareリポジトリを見つけたら
			if isBareRepo(path) {
				fmt.Printf("Found Bare Repo: %s\n", path)
				// TODO: descriptionを抜き取ってbboltへ保存
				return filepath.SkipDir // リポジトリの中までは掘らない
			}
			return nil
		})
	}
	return nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	err = scanRepositories(config)
	if err != nil {
		panic(err)
	}
}
