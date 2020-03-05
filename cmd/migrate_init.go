package cmd

import (
	"fmt"
	"log"

	"github.com/calvinchengx/gin-go-pg/migration"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "creates version info table in the database",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		err := migration.Run(args...)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	migrateCmd.AddCommand(initCmd)
}
