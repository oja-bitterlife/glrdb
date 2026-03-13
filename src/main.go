package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
		lines := strings.SplitSeq(string(out), "\n")
		for line := range lines {
			l := strings.TrimSpace(line)
			if l != "" && !strings.HasPrefix(l, "#") {
				return l // 最初の意味のある行
			}
		}
	}

	return "No description"
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

}
