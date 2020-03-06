package cmd

import (
	"fmt"
	"log"

	"github.com/calvinchengx/gin-go-pg/migration"
	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset all migrations",
	Long:  `reset all migrations`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reset called")
		err := migration.Run("reset")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	migrateCmd.AddCommand(resetCmd)
}
