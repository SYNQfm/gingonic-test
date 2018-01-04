package helpers

import (
	"testing"

	"github.com/SYNQfm/helpers/common"
	"github.com/SYNQfm/helpers/test_helper"
	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	resp := test_helper.Response{}
	assert.NotNil(t, resp)
	err := common.NewError("test error %s", "val")
	assert.Equal(t, "test error val", err.Error())
}
