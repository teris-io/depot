package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

// App represents the depot application configuration.
type Config struct {

	// Hostname to run the depot service on
	Hostname string `json,yaml:"hostname"`

	// Port to run the depot service on
	Port int `json,yaml:"port"`

	// DataDir defines the directory to place artifact repositories
	DataDir string `json,yaml:"dataDir"`

	Remotes []Remote `json,yaml:"remotes"`
}

type Remote struct {
	Name    string
	Url     string
	Headers map[string]string
}

var defaultConfig = Config{
	Hostname: "0.0.0.0",
	Port:     9595,
	DataDir:  "./data",
	Remotes: []Remote{
		{
			Name: "jcenter",
			Url:  "https://jcenter.bintray.com/",
		},
	},
}

// Load loads configuration from a file.
func Load(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

// Parse parses the byte array of configuration data and un-marshals it
// into a Config object.
func Parse(data []byte) (*Config, error) {
	// copy default
	c := defaultConfig
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
