package cmd

import (
	"os"
	"os/exec"
	"reflect"
	"testing"
)

// TestParseFlags ensures given a secret the Config instance is filled
func TestParseFlags(t *testing.T) {
	// secret is required
	secret := []byte("test")
	os.Setenv("VICTIMS_BOT_SECRET", string(secret))
	ParseFlags()
	if Config.Secret != string(secret) {
		t.Errorf("Config.Secret returned bad data. %s != %s", Config.Secret, secret)
		t.Fail()
	}
	foundSecret := Config.GetSecret()
	if !reflect.DeepEqual(foundSecret, secret) {
		t.Errorf("GetSecret returned bad data. %b != %b", foundSecret, secret)
		t.Fail()
	}

	// Use reflection to verify all Config fields have values
	cfg := reflect.ValueOf(Config).Elem()
	for i := 0; i < cfg.NumField(); i++ {
		valueResult := cfg.Field(i)
		typeResult := cfg.Type().Field(i)

		if valueResult.String() == "" {
			t.Errorf("%s should not be empty", typeResult.Name)
		}
	}
}

// TestParseFlagsWithoutSecret ensures that if a secret is not given the bot
// exits with an error
func TestParseFlagsWithoutSecret(t *testing.T) {
	// Ensure secret is empty
	os.Setenv("VICTIMS_BOT_SECRET", "")
	if os.Getenv("EXIT_TEST") == "1" {
		ParseFlags()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestParseFlagsWithoutSecret")
	cmd.Env = append(os.Environ(), "EXIT_TEST=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("ParseFlags should have exited with 1: %v", err)
}
