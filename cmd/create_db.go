package cmd

import (
	"fmt"

	"github.com/calvinchengx/gin-go-pg/config"
	"github.com/calvinchengx/gin-go-pg/manager"
	"github.com/spf13/cobra"
)

// createCmd represents the migrate command
var createdbCmd = &cobra.Command{
	Use:   "create_db",
	Short: "create_db creates a database user and database from database parameters declared in config",
	Long:  `create_db creates a database user and database from database parameters declared in config`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create_db called")
		p := config.GetPostgresConfig()

		// connection to db as postgres superuser
		dbSuper := config.GetPostgresSuperUserConnection()
		defer dbSuper.Close()

		manager.CreateDatabaseUserIfNotExist(dbSuper, p)
		manager.CreateDatabaseIfNotExist(dbSuper, p)
	},
}

func init() {
	rootCmd.AddCommand(createdbCmd)
}
