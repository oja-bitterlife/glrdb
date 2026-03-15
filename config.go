package main

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// **********************************************************************
// Config構造体と関連関数
type Config struct {
	Global  GlobalSection   `toml:"global"`
	Sources []SourceSection `toml:"sources"`
}
type GlobalSection struct {
	DBName   string   `toml:"db_name"`   // データベースファイル名
	MaxDepth int      `toml:"max_depth"` // デフォルトは64、必要に応じて変更可能
	Excludes []string `toml:"exclude"`   // node_modules などをここに追加予定
}

type SourceSection struct {
	Path     string   `toml:"path"`
	Excludes []string `toml:"exclude"` // node_modules などをここに追加予定
}

// **********************************************************************
// Config関連の関数
// ==================================================
// Configのコンストラクタ
func newDefaultGlobalSection() GlobalSection {
	return GlobalSection{
		DBName:   "glrdb.boltdb",
		MaxDepth: 64,
		Excludes: []string{},
	}
}
func newDefaultConfig() *Config {
	return &Config{
		Global: newDefaultGlobalSection(),
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

	return home
}

// **********************************************************************
// Configファイルの読み込み
func loadConfig(path string) (*Config, error) {
	config := newDefaultConfig()
	var data []byte

	// 引数で指定されたパスがあればそこから読み込む
	if path != "" {
		d, err := os.ReadFile(path)
		if err != nil {
			// 引数で指定されたパスが見つからない場合はエラー
			return nil, err
		} else {
			data = d
		}
	} else {
		// 引数がない場合はカレントディレクトリgのglrdb.toml をチェック
		d, err := os.ReadFile(defaultConfigName)
		if err != nil {
			// .config下は今はみない
			return nil, err
		} else {
			data = d
		}
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
