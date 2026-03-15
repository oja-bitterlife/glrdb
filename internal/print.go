package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"go.etcd.io/bbolt"
)

// **********************************************************************
// リポジトリ情報をタブ区切りでbase64で出力する
func printForFzf(config *Config) error {
	db, err := bbolt.Open(config.Global.DBName, 0666, nil) // 読み込みなので 0666
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

			// base64エンコードして出力
			content := base64.StdEncoding.EncodeToString([]byte(repo.Description))

			// タブ区切りで出力（fzfで扱いやすい）
			fmt.Printf("%s\t%s\n", repo.Path, content)
			return nil
		})
	})
}

// **********************************************************************
// リポジトリ情報をsummary形式で出力する
func printList(config *Config) error {
	db, err := bbolt.Open(config.Global.DBName, 0666, nil) // 読み込みなので 0666
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

			summary := repo.Description

			// 改行を<br>に
			summary = strings.ReplaceAll(summary, "\r\n", "\n")
			summary = strings.ReplaceAll(summary, "\r", "\n")
			summary = strings.ReplaceAll(summary, "\n", " ")

			// 連続するスペースを1つのスペースにして、前後のスペースを削除
			strings.Join(strings.Fields(summary), " ")

			// それでも長いときは切る
			if len(summary) > 80 {
				summary = summary[:80] + "..."
			}

			fmt.Printf("%s: %s\n", repo.Path, summary)
			return nil
		})
	})
}
