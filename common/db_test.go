package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDbUrl(t *testing.T) {
	assert := assert.New(t)
	name := parseDatabaseUrl("")
	assert.Equal("", name)
	name = parseDatabaseUrl("://abcd")
	assert.Equal("", name)
	name = parseDatabaseUrl("postgres://user:password@host.com:5432/dbname")
	assert.Equal("host=host.com port=5432 dbname=dbname user=user password=password", name)
	name = parseDatabaseUrl("postgres://user:password@host.com:5432/dbname?sslmode=disable")
	assert.Equal("host=host.com port=5432 dbname=dbname user=user password=password sslmode=disable", name)
	name = parseDatabaseUrl(DEFAULT_DB_URL)
	assert.Equal("host=localhost port=5432 dbname=aerico_test user=circleci password=circleci sslmode=disable", name)
}
