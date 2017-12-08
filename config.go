package main

import (
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"io/ioutil"
)

// AppConfig represents the depot application configuration.
type AppConfig struct {

	// Hostname to run the depot service on
	Hostname string `json,yaml:"hostname"`

	// Port to run the depot service on
	Port int `json,yaml:"port"`
}

// LoadConfig loads configuration from a file.
func LoadConfig(filename string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseConfig(data)
}

// ParseConfig parses the byte array of configuration data and un-marshals it
// into an AppConfig object.
func ParseConfig(data []byte) (*AppConfig, error) {
	c := defaultAppConfig()
	var err error
	if err = yaml.Unmarshal(data, c); err != nil {
		return nil, errors.Wrap(err, "yaml unmarshaling failed")
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func defaultAppConfig() *AppConfig {
	return &AppConfig{
		Hostname: "0.0.0.0",
		Port:     9595,
	}
}
