package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// TwilioConfig persists the config for our Twilio services
type TwilioConfig struct {
	Account    string `env:"TWILIO_ACCOUNT"`
	Token      string `env:"TWILIO_TOKEN"`
	VerifyName string `env:"TWILIO_VERIFY_NAME"`
	Verify     string `env:"TWILIO_VERIFY"`
}

// GetTwilioConfig returns a TwilioConfig pointer with the correct Mail Config values
func GetTwilioConfig() *TwilioConfig {
	c := TwilioConfig{}
	if err := env.Parse(&c); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &c
}
