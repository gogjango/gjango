package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gogjango/gjango/route"
	"github.com/gogjango/gjango/server"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// routes will be attached to s
var s server.Server

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gjango",
	Short: "A simple golang framework for API server",
	Long: `gjango is a simple golang framework for building high performance, easy-to-extend API web servers.
	Inspired by django python web framework, gjango aims to make it simple and fast to build production grade, web applications.
	
	By default, our program will run the API server.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var env string
		var ok bool
		if env, ok = os.LookupEnv("GJANGO_ENV"); !ok {
			env = "dev"
			fmt.Printf("Run server in %s mode\n", env)
		}
		fmt.Println(s)
		err := s.Run(env)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(customRouteServices []route.ServicesI) {
	s.RouteServices = customRouteServices
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gjango.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gjango" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gjango")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
