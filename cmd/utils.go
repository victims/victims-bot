package cmd

import (
	"flag"
	"os"
)

// Configuration holds all configuration parameters
type Configuration struct {
	Secret string
	BindTo string
}

// Config is the singleton Configuration instance
var Config *Configuration

func init() {
	if Config == nil {
		Config = &Configuration{}
	}
}

// GetSecret returns the secret in byte form
func (c *Configuration) GetSecret() []byte {
	return []byte(c.Secret)
}

// ParseFlags parses CLI flags and enforces requirements
func ParseFlags() {
	flag.StringVar(&Config.BindTo, "bind", "0.0.0.0:9999", "Host:Port to bind on")
	flag.StringVar(&Config.Secret, "secret", "", "Shared secret with github")
	flag.Parse()

	// secret must be given
	if Config.Secret == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
