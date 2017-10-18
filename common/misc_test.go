package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseType(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("errored", ParseType("Error"))
	assert.Equal("count", ParseType("ct"))
	assert.Equal("skipped", ParseType("SKiP"))
	assert.Equal("already", ParseType("already"))
}

func TestRet(t *testing.T) {
	assert := assert.New(t)
	ret := NewRet("test")
	assert.Equal(0, ret.Value("count"))
	ret.Add("count")
	assert.True(ret.Eq("count", 1))
	assert.True(ret.Gte("count", 1))
	assert.True(ret.Lte("count", 1))
	assert.False(ret.Lt("count", 1))
	assert.False(ret.Gt("count", 1))
}
