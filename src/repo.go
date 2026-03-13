package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func fetchReadme(path string) string {
	// descriptionファイルがあれば優先
	if data, err := os.ReadFile(filepath.Join(path, "description")); err == nil {
		desc := strings.TrimSpace(string(data))
		// Gitデフォルトの文言でなければ採用
		if desc != "" && !strings.Contains(desc, "Unnamed repository") {
			return desc
		}
	}

	// いくつかの候補から最初に見つかったものを採用
	targets := []string{"README.md", "Readme.md", "readme.md"}
	for _, t := range targets {
		cmd := exec.Command("git", "-C", path, "show", "HEAD:"+t)
		out, err := cmd.Output()
		if err == nil {
			return string(out)
		}
	}

	return ""
}
