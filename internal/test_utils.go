// Package testutils provides internal shared testing utilities.
package testutils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	mathrand "math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/jswny/sad"
)

// StringOptions represents all options as strings.
type StringOptions struct {
	Server     string
	Username   string
	RootDir    string
	PrivateKey string
	Channel    string
	Path       string
	EnvVars    string
	Debug      string
}

// FromOptions converts options into string options.
func (stringOpts *StringOptions) FromOptions(opts *sad.Options) {
	stringOpts.Server = opts.Server.String()
	stringOpts.Username = opts.Username
	stringOpts.RootDir = opts.RootDir
	stringOpts.PrivateKey = opts.PrivateKey.ToBase64PEMString()
	stringOpts.Channel = opts.Channel
	stringOpts.Path = opts.Path
	stringOpts.EnvVars = strings.Join(opts.EnvVars, ",")
	stringOpts.Debug = strconv.FormatBool(opts.Debug)
}

// SetEnv sets environment variables for all string options.
// UnsetEnv should be called after.
func (stringOpts *StringOptions) SetEnv() {
	prefix := sad.EnvVarPrefix
	var envVarPostfix string

	envVarPostfix = "SERVER"
	os.Setenv(prefix+envVarPostfix, stringOpts.Server)

	envVarPostfix = "USERNAME"
	os.Setenv(prefix+envVarPostfix, stringOpts.Username)

	envVarPostfix = "ROOT_DIR"
	os.Setenv(prefix+envVarPostfix, stringOpts.RootDir)

	envVarPostfix = "PRIVATE_KEY"
	os.Setenv(prefix+envVarPostfix, stringOpts.PrivateKey)

	envVarPostfix = "CHANNEL"
	os.Setenv(prefix+envVarPostfix, stringOpts.Channel)

	envVarPostfix = "PATH"
	os.Setenv(prefix+envVarPostfix, stringOpts.Path)

	envVarPostfix = "ENV_VARS"
	os.Setenv(prefix+envVarPostfix, stringOpts.EnvVars)

	envVarPostfix = "DEBUG"
	os.Setenv(prefix+envVarPostfix, stringOpts.Debug)
}

// UnsetEnv sets environment variables for all string options.
// Should be called after SetEnv.
func (stringOpts *StringOptions) UnsetEnv() {
	prefix := sad.EnvVarPrefix
	var envVarPostfix string

	envVarPostfix = "SERVER"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "USERNAME"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "ROOT_DIR"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "PRIVATE_KEY"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "CHANNEL"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "PATH"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "ENV_VARS"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "DEBUG"
	defer os.Unsetenv(prefix + envVarPostfix)
}

// GetTestOpts retrieves a set of random options for testing.
func GetTestOpts() sad.Options {
	rsaPrivateKey := GenerateRSAPrivateKey()

	randSize := 5

	testOpts := sad.Options{
		Server:     net.ParseIP("1.2.3.4"),
		Username:   randString(randSize),
		RootDir:    randString(randSize),
		PrivateKey: rsaPrivateKey,
		Channel:    randString(randSize),
		Path:       randString(randSize),
		EnvVars: []string{
			randString(randSize),
			randString(randSize),
		},
		Debug: true,
	}

	return testOpts
}

// GenerateRSAPrivateKey generates a random RSA private key.
func GenerateRSAPrivateKey() sad.RSAPrivateKey {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	rsaPrivateKey := sad.RSAPrivateKey{
		PrivateKey: privateKey,
	}
	return rsaPrivateKey
}

// CompareOpts compares two sets of options in a test environment.
func CompareOpts(expectedOpts sad.Options, actualOpts sad.Options, t *testing.T) {
	if !actualOpts.Server.Equal(expectedOpts.Server) {
		t.Errorf("Expected server IP %s but got %s", expectedOpts.Server, actualOpts.Server)
	}

	compareStrings(expectedOpts.Username, actualOpts.Username, "username", t)

	compareStrings(expectedOpts.RootDir, actualOpts.RootDir, "root directory", t)

	if (expectedOpts.PrivateKey.PrivateKey != nil && actualOpts.PrivateKey.PrivateKey != nil) && !expectedOpts.PrivateKey.PrivateKey.Equal(actualOpts.PrivateKey.PrivateKey) {
		t.Errorf("Expected equal private keys but they were not equal")
	}

	compareStrings(expectedOpts.Channel, actualOpts.Channel, "channel", t)

	if actualOpts.Path != expectedOpts.Path {
		t.Errorf("Expected path %s but got %s", expectedOpts.Path, actualOpts.Path)
	}
	compareStrings(expectedOpts.Path, actualOpts.Path, "path", t)

	compareSlices(actualOpts.EnvVars, expectedOpts.EnvVars, "environment variables", t)

	if actualOpts.Debug != expectedOpts.Debug {
		t.Errorf("Expected debug %t but got %t", expectedOpts.Debug, actualOpts.Debug)
	}
}

func compareStrings(expected, actual, name string, t *testing.T) {
	if expected != actual {
		empty := "<empty>"

		if expected == "" {
			expected = empty
		}

		if actual == "" {
			actual = empty
		}

		t.Errorf("Expected %s %s but got %s", name, expected, actual)
	}
}

// CloneOptions clones options into other options.
// The options to clone into should ideally be empty.
func CloneOptions(optionsToClone *sad.Options, optionsToCloneInto *sad.Options) error {
	data, err := json.Marshal(optionsToClone)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &optionsToCloneInto)
	if err != nil {
		return err
	}

	return nil
}

func randString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(b)
}

func compareSlices(expected, actual []string, name string, t *testing.T) {
	equal := true

	if (expected == nil) != (actual == nil) {
		equal = false
	}

	if len(expected) != len(actual) {
		equal = false
	}

	for i := range expected {
		if expected[i] != actual[i] {
			equal = false
		}
	}

	if !equal {
		t.Errorf("Expected %s %s but got %s", name, expected, actual)
	}
}
