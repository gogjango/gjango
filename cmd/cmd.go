package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/calvinchengx/gin-go-pg/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "run runs our API server",
	Long:  `go run . run gets our API server running`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		err := server.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute is the entry point for all our commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
