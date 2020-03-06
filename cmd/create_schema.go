package cmd

import (
	"fmt"

	"github.com/calvinchengx/gin-go-pg/config"
	"github.com/calvinchengx/gin-go-pg/manager"
	"github.com/calvinchengx/gin-go-pg/repository"
	"github.com/calvinchengx/gin-go-pg/secret"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// createschemaCmd represents the createschema command
var createSchemaCmd = &cobra.Command{
	Use:   "create_schema",
	Short: "create_schema creates the initial database schema for the existing database",
	Long:  `create_schema creates the initial database schema for the existing database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createschema called")

		db := config.GetConnection()
		log, _ := zap.NewDevelopment()
		defer log.Sync()
		accountRepo := repository.NewAccountRepo(db, log, secret.New())
		roleRepo := repository.NewRoleRepo(db, log)

		m := manager.NewManager(accountRepo, roleRepo, db)
		models := manager.GetModels()
		m.CreateSchema(models...)
		m.CreateRoles()
	},
}

func init() {
	rootCmd.AddCommand(createSchemaCmd)
}
