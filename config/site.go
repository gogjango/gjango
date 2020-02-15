package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// SiteConfig persists global configs needed for our application
type SiteConfig struct {
	ExternalURL string `env:"EXTERNAL_URL"  envDefault:"http://localhost:8080"`
}

// GetSiteConfig returns a SiteConfig pointer with the correct Site Config values
func GetSiteConfig() *SiteConfig {
	c := SiteConfig{}
	if err := env.Parse(&c); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &c
}
