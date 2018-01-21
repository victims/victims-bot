package cmd

import (
	"flag"
	"os"
)

// Configuration holds all configuration parameters
type Configuration struct {
	Secret         string
	BindTo         string
	GitHubUsername string
	GitHubPassword string
	GitRepo        string
}

// Config is the singleton Configuration instance
var Config *Configuration

// init ensures the creation of the Config singleton
func init() {
	if Config == nil {
		Config = &Configuration{}
	}
}

// GetSecret returns the secret in byte form
func (c *Configuration) GetSecret() []byte {
	return []byte(c.Secret)
}

// getEnv returns either the environment variable or the default
func getEnv(variable, defaultValue string) string {
	result := os.Getenv(variable)
	if result == "" {
		return defaultValue
	}
	return result
}

// ParseFlags parses CLI flags and enforces requirements
func ParseFlags() {
	flag.StringVar(&Config.BindTo, "bind",
		getEnv("VICTIMS_BOT_BIND", "0.0.0.0:9999"), "Host:Port to bind on")
	flag.StringVar(&Config.Secret, "secret",
		getEnv("VICTIMS_BOT_SECRET", ""), "Shared secret with github")
	flag.StringVar(&Config.GitHubUsername, "github-username",
		getEnv("VICTIMS_BOT_GITHUB_USERNAME", "victims-bot"),
		"Name of the bot in GitHub to know what events to ignore")
	flag.StringVar(&Config.GitHubPassword, "github-password",
		getEnv("VICTIMS_BOT_GITHUB_PASSWORD", ""),
		"Password used to access GitHub")
	flag.StringVar(
		&Config.GitRepo, "git-repo",
		getEnv("VICTIMS_BOT_GIT_REPO", "git@github.com:victims/victims-cve-db.git"),
		"git repo to clone and push using ssh")

	// Parse CLI and store values
	flag.Parse()

	// secret must be given
	if Config.Secret == "" || Config.GitHubPassword == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
