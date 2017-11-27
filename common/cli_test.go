package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCli(t *testing.T) {
	assert := assert.New(t)
	cli := NewCli()
	cli.String("command", "test", "this is the test")
	cli.Parse([]string{"-command", "cmd"})
	assert.Equal("cmd", cli.GetString("command"))
}

func TestDefaultCli(t *testing.T) {
	assert := assert.New(t)
	cli := NewCli()
	cli.DefaultSetup("test", "this is my command")
	cli.Parse([]string{})
	assert.Equal("test", cli.GetString("command"))
	assert.Equal(120, cli.GetInt("timeout"))
	assert.True(cli.Simulate)
	assert.Equal(time.Duration(120)*time.Second, cli.Timeout)
	assert.Equal(10, cli.GetInt("limit"))
}
