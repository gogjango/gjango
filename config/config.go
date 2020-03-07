package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Load returns Configuration struct
func Load(env string) *Configuration {
	_, filePath, _, _ := runtime.Caller(0)
	configName := "config." + env + ".yaml"
	configPath := filePath[:len(filePath)-9] + "files" + string(filepath.Separator)

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var config Configuration
	viper.Unmarshal(&config)

	return &config
}

// Configuration holds data necessery for configuring application
type Configuration struct {
	Server *Server `yaml:"server"`
	JWT    *JWT    `yaml:"jwt"`
}

// Server holds data necessary for server configuration
type Server struct {
	Port int `yaml:"port"`
}

// JWT holds data necessary for JWT configuration
type JWT struct {
	Realm            string
	Secret           string
	Duration         int
	RefreshDuration  int
	MaxRefresh       int
	SigningAlgorithm string
}
