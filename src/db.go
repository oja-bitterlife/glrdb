package main

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

func updateDB(config *Config, repos []string) error {
	db, err := bbolt.Open(config.DBName, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Print("start transaxtion\n")
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
				Path:        path,
				Description: desc,
			}

			data, _ := json.Marshal(repo)
			if err := bucket.Put([]byte(path), data); err != nil {
				return err
			}
		}
		return nil
	})
}
