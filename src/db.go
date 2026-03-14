package main

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

type Repository struct {
	Path        string `json:"path"`
	IsBare      bool   `json:"is_bare"`
	Description string `json:"description"`
}

func updateDB(config *Config, repos []repoPath) error {
	db, err := bbolt.Open(config.Global.DBName, 0644, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bbolt.Tx) error {
		// バケット作成
		bucket, err := tx.CreateBucketIfNotExists([]byte("Repositories"))
		if err != nil {
			return err
		}

		for _, path := range repos {
			// 前に作った fetchReadme を使用
			desc := fetchReadme(path)
			repo := Repository{
				Path:        path.path,
				IsBare:      path.isBare,
				Description: desc,
			}

			data, _ := json.Marshal(repo)
			if err := bucket.Put([]byte(path.path), data); err != nil {
				return err
			}
		}
		return nil
	})
}
