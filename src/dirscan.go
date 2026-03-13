package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func isBlacklisted(name string, blacklist []string) bool {
	return slices.Contains(blacklist, name)
}

func isBareRepo(path string) bool {
	_, errObjects := os.Stat(filepath.Join(path, "objects"))
	_, errRefs := os.Stat(filepath.Join(path, "refs"))
	_, errConfig := os.Stat(filepath.Join(path, "config"))
	_, errHead := os.Stat(filepath.Join(path, "HEAD"))
	return errObjects == nil && errRefs == nil && errConfig == nil && errHead == nil
}

func scanDir(config *Config) ([]string, error) {
	var allRepos []string

	for _, src := range config.Sources {
		fmt.Printf("Scanning: %s\n", src.Path)
		repos, err := rec_scanDir(src.Path, 0, config)
		if err != nil {
			return nil, err
		}
		if repos != nil {
			allRepos = append(allRepos, repos...)
		}
	}

	// debug
	fmt.Printf("Total repositories found: %#v\n", allRepos)

	return allRepos, nil
}

func rec_scanDir(path string, depth int, config *Config) ([]string, error) {
	var foundRepos []string

	if depth > config.MaxDepth {
		return nil, fmt.Errorf("max depth exceeded at: %s", path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, nil // アクセス権限エラーなどは無視
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// シンボリックリンクの解決
		if info.Mode()&os.ModeSymlink != 0 {
			resolved, err := filepath.EvalSymlinks(fullPath)
			if err != nil {
				continue
			}
			fullPath = resolved
			info, _ = os.Stat(fullPath)
		}

		if info.IsDir() {
			name := entry.Name()

			// スキップ判定（ドット、blacklist）
			if strings.HasPrefix(name, ".") || isBlacklisted(name, config.Blacklist) {
				continue
			}

			// リポジトリ判定
			if isBareRepo(fullPath) {
				foundRepos = append(foundRepos, fullPath)
				continue // リポジトリの中は掘らない
			}

			// 再帰呼び出しの結果をマージする
			subRepos, err := rec_scanDir(fullPath, depth+1, config)
			if err != nil {
				return nil, err
			}
			foundRepos = append(foundRepos, subRepos...)
		}
	}
	return nil
}
