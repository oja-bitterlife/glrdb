package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// バージョンは頻繁に書き換えるので個別で定数化しておく
const version = "0.1.0"

// Default定数
const (
	// コンフィグのファイル名
	defaultConfigName = "glrdb.toml"
)

func main() {
	app := &cli.App{
		Name:    "glrdb",
		Usage:   "A tool to manage your git local repositories with descriptions",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   defaultConfigName,
				Aliases: []string{"f"},
				Usage:   "Path to the config file (default: glrdb.toml)",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "update",
				Usage: "Scan directories and update the database with repository descriptions",
				Action: func(ctx *cli.Context) error {
					if config, err := loadConfig(ctx.String("config")); err != nil {
						return err
					} else {
						fmt.Printf("--- start scanning repositories ---\n")
						allRepos, err := scanDir(config)
						if err != nil {
							return err
						}

						fmt.Printf("\n--- start fetching description ---\n")
						if err = updateDB(config, allRepos); err != nil {
							return err
						}
						return nil
					}
				},
			},
			{
				Name:  "print",
				Usage: "Print repository information in a format suitable for fzf",
				Action: func(ctx *cli.Context) error {
					if config, err := loadConfig(ctx.String("config")); err != nil {
						return err
					} else {
						printForFzf(config)
						return nil
					}
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
