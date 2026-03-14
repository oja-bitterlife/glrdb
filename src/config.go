package main

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// **********************************************************************
// Config構造体と関連関数
type Config struct {
	DBName    string   `toml:"db_name"` // データベースファイル名
	Sources   []Source `toml:"sources"`
	Blacklist []string `toml:"blacklist"` // node_modules などをここに追加予定
	MaxDepth  int      `toml:"max_depth"` // デフォルトは64、必要に応じて変更可能
}

type Source struct {
	Path string `toml:"path"`
}

// **********************************************************************
// Config関連の関数
// ==================================================
// Configのコンストラクタ
func newDefaultConfig() *Config {
	return &Config{
		DBName:   "glrdb.boltdb",
		MaxDepth: 64,
		// とりあえずデフォルトを手入力
		Blacklist: []string{"node_modules", "vendor"},
	}
}

// ==================================================
// ~/ をホームディレクトリに展開する
func expandHome(path string) string {
	// チルダで始まっていない場合はそのまま返す
	if path == "" || path[0] != '~' {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path // ホームディレクトリが取得できない場合はそのまま返す
	}

	// "~/path/to/dir" の場合は "~/ "の2文字を home で置き換える
	if len(path) > 1 && path[1] == '/' {
		return filepath.Join(home, path[2:])
	}

	// "~" 単体の場合
	return home
}

// **********************************************************************
// Configファイルの読み込み
func loadConfig(path string) (*Config, error) {
	config := newDefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	// パスを展開しておく
	for i, src := range config.Sources {
		config.Sources[i].Path = expandHome(src.Path)
	}

	return config, nil
}
