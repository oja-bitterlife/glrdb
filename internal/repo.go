package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// **********************************************************************
func findMesage(path repoPath, hitFile string) {
	fmt.Printf("%s: %s\n", path.path, hitFile)
}

// リポジトリのdescriptionかREADME.mdを取得する
func fetchReadme(path repoPath) string {
	var gitPathArg []string

	if path.isBare {
		// descriptionファイルがあれば優先
		if data, err := os.ReadFile(filepath.Join(path.path, "description")); err == nil {
			desc := strings.TrimSpace(string(data))
			// Gitデフォルトの文言でなければ採用
			if desc != "" && !strings.Contains(desc, "Unnamed repository") {
				findMesage(path, "description")
				return desc
			}
		}

		// なければreadmeを探すためにgitコマンドでアクセス
		gitPathArg = []string{"-C", path.path}
	} else {
		// 通常のリポジトリは.gitディレクトリにアクセスする
		gitPathArg = []string{"--git-dir", filepath.Join(path.path, ".git")}
	}

	// HEAD のツリーからファイル名一覧を取得
	gitArgs := append(gitPathArg, "ls-tree", "--name-only", "HEAD")
	cmd := exec.Command("git", gitArgs...)
	out, err := cmd.Output()
	if err != nil {
		return "" // HEADがない（空の）リポジトリなど
	}
	files := strings.Split(string(out), "\n")

	// 優先順位を決めてファイルを探す
	var targetFile string
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
	if targetFile != "" {
		findMesage(path, targetFile)

		gitArgs := append(gitPathArg, "show", "HEAD:"+targetFile)
		cmd := exec.Command("git", gitArgs...)
		if out, err = cmd.Output(); err == nil {
			return string(out)
		}
	} else {
		fmt.Printf("%s:\n", path.path)
	}

	return ""
}
