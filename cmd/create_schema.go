package cmd

import (
	"fmt"

	"github.com/calvinchengx/gin-go-pg/config"
	"github.com/calvinchengx/gin-go-pg/manager"
	"github.com/spf13/cobra"
)

// createschemaCmd represents the createschema command
var createSchemaCmd = &cobra.Command{
	Use:   "create_schema",
	Short: "create_schema creates the initial database schema for the existing database",
	Long:  `create_schema creates the initial database schema for the existing database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createschema called")
		db := config.GetConnection()
		models := manager.GetModels()
		manager.CreateSchema(db, models...)
	},
}

func init() {
	rootCmd.AddCommand(createSchemaCmd)
}
