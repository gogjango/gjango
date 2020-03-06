package cmd

import (
	"fmt"

	"github.com/calvinchengx/gin-go-pg/config"
	"github.com/calvinchengx/gin-go-pg/manager"
	"github.com/spf13/cobra"
)

// createCmd represents the migrate command
var createdbCmd = &cobra.Command{
	Use:   "createdb",
	Short: "createdb creates a database from database parameters declared in config",
	Long:  `createdb creates a database from database parameters declared in config`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createdb called")
		p := config.GetPostgresConfig()

		db := config.GetConnection()
		defer db.Close()

		manager.CreateDatabaseIfNotExist(db, p)
	},
}

func init() {
	rootCmd.AddCommand(createdbCmd)
}
