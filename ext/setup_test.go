package ext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSetup(t *testing.T) {
	assert := assert.New(t)
	setup := GetSetupByEnv("")
	assert.Equal("", setup.Version)
	setup = GetSetupByEnv("v2")
	assert.Equal("v2", setup.Version)
}

func TestSetupSynq(t *testing.T) {
	assert := assert.New(t)
	api := SetupSynq()
	assert.Equal("v1", api.Version())
}
