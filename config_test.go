package main_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/teris-io/depot"
	"testing"
)

func TestParse_empty_yieldsDefault(t *testing.T) {
	str := ""
	c, err := main.ParseConfig([]byte(str))
	assert.Nil(t, err)
	assert.Equal(t, "0.0.0.0", c.Hostname)
	assert.Equal(t, 9595, c.Port)
}

func TestParse_hostPort_configurable(t *testing.T) {
	str := "hostname: 0.1.2.3\nport: 1234"
	c, err := main.ParseConfig([]byte(str))
	assert.Nil(t, err)
	assert.Equal(t, "0.1.2.3", c.Hostname)
	assert.Equal(t, 1234, c.Port)
}

func TestParse_partiallyMissing_yieldsDefault(t *testing.T) {
	str := "port: 1234"
	c, err := main.ParseConfig([]byte(str))
	assert.Nil(t, err)
	assert.Equal(t, "0.0.0.0", c.Hostname)
	assert.Equal(t, 1234, c.Port)
}
