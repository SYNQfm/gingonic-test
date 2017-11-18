package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDbUrl(t *testing.T) {
	assert := assert.New(t)
	name := ParseDatabaseUrl("")
	assert.Equal("", name)
	name = ParseDatabaseUrl("://abcd")
	assert.Equal("", name)
	name = ParseDatabaseUrl("postgres://user:password@host.com:5432/dbname")
	assert.Equal("host=host.com port=5432 dbname=dbname user=user password=password", name)
	name = ParseDatabaseUrl("postgres://user:password@host.com:5432/dbname?sslmode=disable")
	assert.Equal("host=host.com port=5432 dbname=dbname user=user password=password sslmode=disable", name)
	name = ParseDatabaseUrl(DEFAULT_DB_URL)
	assert.Equal("host=localhost port=5432 dbname=db_test user=circleci password=circleci sslmode=disable", name)
}
