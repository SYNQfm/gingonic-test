package helpers

import (
	"testing"

	"github.com/SYNQfm/test_helper"
	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	assert.NotNil(t, test_helper.VIDEO_ID)
}
