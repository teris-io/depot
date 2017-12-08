package config_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/teris-io/depot/config"
	"testing"
)

func TestParse_empty_yieldsDefault(t *testing.T) {
	str := ""
	c, err := config.Parse([]byte(str))
	assert.Nil(t, err)
	assert.Equal(t, "0.0.0.0", c.Hostname)
	assert.Equal(t, 9595, c.Port)
}

func TestParse_hostPort_configurable(t *testing.T) {
	str := "hostname: 0.1.2.3\nport: 1234\ndataDir: \"some dir\""
	c, err := config.Parse([]byte(str))
	assert.Nil(t, err)
	assert.Equal(t, "0.1.2.3", c.Hostname)
	assert.Equal(t, 1234, c.Port)
	assert.Equal(t, "some dir", c.DataDir)
}

func TestParse_partiallyMissing_yieldsDefault(t *testing.T) {
	str := "port: 1234"
	c, err := config.Parse([]byte(str))
	assert.Nil(t, err)
	assert.Equal(t, "0.0.0.0", c.Hostname)
	assert.Equal(t, 1234, c.Port)
}
