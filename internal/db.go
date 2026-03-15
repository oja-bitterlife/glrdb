package main

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

// **********************************************************************
// データベース関連の構造体
type Repository struct {
	Path        string `json:"path"`
	IsBare      bool   `json:"is_bare"`
	Description string `json:"description"`
}

// **********************************************************************
// データベースの更新
func updateDB(config *Config, repos []repoPath) error {
	// データベースを開く（存在しない場合は新規作成）
	db, err := bbolt.Open(config.Global.DBName, 0644, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	// トランザクションを開始してデータベースを更新
	return db.Update(func(tx *bbolt.Tx) error {
		// いったん全クリア
		tx.DeleteBucket([]byte("Repositories"))

		// バケット作成
		bucket, err := tx.CreateBucketIfNotExists([]byte("Repositories"))
		if err != nil {
			return err
		}

		// リポジトリごとにデータを保存
		for _, path := range repos {
			repo := Repository{
				Path:        path.path,
				IsBare:      path.isBare,
				Description: fetchReadme(path),
			}

			if data, err := json.Marshal(repo); err != nil {
				return err
			} else {
				if err := bucket.Put([]byte(path.path), data); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
