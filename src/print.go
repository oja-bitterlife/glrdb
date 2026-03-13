package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.etcd.io/bbolt"
)

func printForFzf(config *Config) error {
	db, err := bbolt.Open(config.DBName, 0666, nil) // 読み込みなので 0666
	if err != nil {
		return err
	}
	defer db.Close()

	return db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("Repositories"))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(key, value []byte) error {
			var repo Repository
			json.Unmarshal(value, &repo)

			// READMEの中身を改行除去して1行に短縮（fzfで見やすくするため）
			summary := strings.ReplaceAll(repo.Description, "\n", " ")
			// if len(summary) > 100 {
			// 	summary = summary[:100] + "..."
			// }

			// タブ区切りで出力（fzfで扱いやすい）
			fmt.Printf("%s\t%s\t%s\n", repo.Name, summary, repo.Path)
			return nil
		})
	})
}
