package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// MailConfig persists the config for our PostgreSQL database connection
type MailConfig struct {
	Name  string `env:"DEFAULT_NAME"`
	Email string `env:"DEFAULT_EMAIL"`
}

// GetMailConfig returns a MailConfig pointer with the correct Mail Config values
func GetMailConfig() *MailConfig {
	c := MailConfig{}
	if err := env.Parse(&c); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &c
}
