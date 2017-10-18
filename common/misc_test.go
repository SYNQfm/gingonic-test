package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRet(t *testing.T) {
	assert := assert.New(t)
	ret := NewRet("test")
	assert.Equal(0, ret.Value("count"))
}
