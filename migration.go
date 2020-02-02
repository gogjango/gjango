package main

import (
	"log"
	"os"
	"path"

	"github.com/calvinchengx/gin-go-pg/config"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

const directory = "migrations"

// handle schema migrations, invoke with `go run . [command]`
func handleMigration(args []string) {
	if len(args) >= 2 {

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		migrationPath := path.Join(cwd, directory)

		if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
			os.Mkdir(migrationPath, 0755)
		}

		// migrations
		db := config.GetConnection()
		defer db.Close()
		err = migrations.Run(db, directory, os.Args)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}
}
