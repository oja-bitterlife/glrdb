package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Repository struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}

func fetchReadme(path string) string {
	// descriptionファイルがあれば優先
	if data, err := os.ReadFile(filepath.Join(path, "description")); err == nil {
		desc := strings.TrimSpace(string(data))
		// Gitデフォルトの文言でなければ採用
		if desc != "" && !strings.Contains(desc, "Unnamed repository") {
			fmt.Printf("Found description file in %s\n", path)
			return desc
		}
	}

	// HEAD のツリーからファイル名一覧を取得
	cmd := exec.Command("git", "-C", path, "ls-tree", "--name-only", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "" // HEADがない（空の）リポジトリなど
	}

	files := strings.Split(string(out), "\n")
	var targetFile string

	// 優先順位を決めてファイルを探す
	for _, file := range files {
		lowName := strings.ToLower(strings.TrimSpace(file))
		if strings.HasPrefix(lowName, "readme") {
			targetFile = file
			// .md を見つけたら即確定、そうでなければとりあえず保持して続行
			if strings.HasSuffix(lowName, ".md") {
				break
			}
		}
	}

	// 特定したファイル名で中身を取得
	fmt.Printf("Looking for README in %s, target: %s\n", path, targetFile)
	if targetFile != "" {
		cmd = exec.Command("git", "-C", path, "show", "HEAD:"+targetFile)
		if out, err = cmd.Output(); err == nil {
			return string(out)
		}
	}

	return ""
}
