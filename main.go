package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"go.etcd.io/bbolt"
)

// バージョンは頻繁に書き換えるので個別で定数化しておく
const version = "0.1.0"

const (
	// データベースのファイル名
	dbFileName = "glrdb.boltdb"
	// コンフィグのファイル名
	configFileName = "glrdb.toml"
)

type Repository struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Bareリポジトリから説明文を取得する
func fetchDescription(repoPath string) string {
	// 1. descriptionファイルがあれば優先
	if data, err := os.ReadFile(filepath.Join(repoPath, "description")); err == nil {
		d := strings.TrimSpace(string(data))
		// Gitデフォルトの文言でなければ採用
		if d != "" && !strings.Contains(d, "Unnamed repository") {
			return d
		}
	}

	// 2. なければ README.md の1行目を git show で試みる
	// Bareなので直接ファイルは読めないため git コマンドを使う
	cmd := exec.Command("git", "-C", repoPath, "show", "HEAD:README.md")
	if out, err := cmd.Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			l := strings.TrimSpace(line)
			if l != "" && !strings.HasPrefix(l, "#") {
				return l // 最初の意味のある行
			}
		}
	}

	return "No description"
}

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
func scanRepositories(db *bbolt.DB, config *Config) error {
	return db.Update(func(tx *bbolt.Tx) error {
		// バケット（テーブルのようなもの）の準備
		b, err := tx.CreateBucketIfNotExists([]byte("repos"))
		if err != nil {
			return err
		}

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

					repo := Repository{
						Path: path,
						Name: name,
						// description取得は一旦スタブ（仮）でもOK
						Description: "Sample description",
					}

					// JSONにして保存
					v, _ := json.Marshal(repo)
					b.Put([]byte(path), v)

					return filepath.SkipDir // リポジトリの中までは掘らない
				}
				return nil
			})
		}
		return nil
	})
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	db, err := bbolt.Open(dbFileName, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = scanRepositories(db, config)
	if err != nil {
		panic(err)
	}
}
